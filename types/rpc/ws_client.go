package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rcrowley/go-metrics"
	tmjson "github.com/tendermint/tendermint/libs/json"
	tmpubsub "github.com/tendermint/tendermint/libs/pubsub"

	tmrand "github.com/tendermint/tendermint/libs/rand"
	"github.com/tendermint/tendermint/libs/service"
	tmsync "github.com/tendermint/tendermint/libs/sync"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/rpc/jsonrpc/types"
)

const (
	defaultMaxReconnectAttempts = 25
	defaultWriteWait            = 0
	defaultReadWait             = 0
	defaultPingPeriod           = 0
)

const (
	protoHTTP  = "http"
	protoHTTPS = "https"
	protoWSS   = "wss"
	protoWS    = "ws"
	protoTCP   = "tcp"
)

//-----------------------------------------------------------------------------
// WSEvents

// WSEvents is a wrapper around WSClient, which implements EventsClient.
type WSEvents struct {
	service.BaseService
	remote   string
	endpoint string
	header   http.Header
	ws       *WSClient

	mtx           tmsync.RWMutex
	subscriptions map[string]chan ctypes.ResultEvent // query -> chan
}

func newWSEvents(remote string, endpoint string, header http.Header) (*WSEvents, error) {
	w := &WSEvents{
		endpoint:      endpoint,
		remote:        remote,
		header:        header,
		subscriptions: make(map[string]chan ctypes.ResultEvent),
	}
	w.BaseService = *service.NewBaseService(nil, "WSEvents", w)

	var err error
	w.ws, err = NewWS(w.remote, w.endpoint, w.header, OnReconnect(func() {
		// resubscribe immediately
		w.redoSubscriptionsAfter(0 * time.Second)
	}))
	if err != nil {
		return nil, err
	}
	w.ws.SetLogger(w.Logger)
	return w, nil
}

// OnStart implements service.Service by starting WSClient and event loop.
func (w *WSEvents) OnStart() error {
	if err := w.ws.Start(); err != nil {
		return err
	}

	go w.eventListener()

	return nil
}

// OnStop implements service.Service by stopping WSClient.
func (w *WSEvents) OnStop() {
	if err := w.ws.Stop(); err != nil {
		w.Logger.Error("Can't stop ws client", "err", err)
	}
}

// Subscribe implements EventsClient by using WSClient to subscribe given
// subscriber to query. By default, returns a channel with cap=1. Error is
// returned if it fails to subscribe.
//
// Channel is never closed to prevent clients from seeing an erroneous event.
//
// It returns an error if WSEvents is not running.
func (w *WSEvents) Subscribe(ctx context.Context, subscriber, query string,
	outCapacity ...int) (out <-chan ctypes.ResultEvent, err error) {

	if !w.IsRunning() {
		return nil, errNotRunning
	}

	if err := w.ws.Subscribe(ctx, query); err != nil {
		return nil, err
	}

	outCap := 1
	if len(outCapacity) > 0 {
		outCap = outCapacity[0]
	}

	outc := make(chan ctypes.ResultEvent, outCap)
	w.mtx.Lock()
	// subscriber param is ignored because Tendermint will override it with
	// remote IP anyway.
	w.subscriptions[query] = outc
	w.mtx.Unlock()

	return outc, nil
}

// Unsubscribe implements EventsClient by using WSClient to unsubscribe given
// subscriber from query.
//
// It returns an error if WSEvents is not running.
func (w *WSEvents) Unsubscribe(ctx context.Context, subscriber, query string) error {
	if !w.IsRunning() {
		return errNotRunning
	}

	if err := w.ws.Unsubscribe(ctx, query); err != nil {
		return err
	}

	w.mtx.Lock()
	_, ok := w.subscriptions[query]
	if ok {
		delete(w.subscriptions, query)
	}
	w.mtx.Unlock()

	return nil
}

// UnsubscribeAll implements EventsClient by using WSClient to unsubscribe
// given subscriber from all the queries.
//
// It returns an error if WSEvents is not running.
func (w *WSEvents) UnsubscribeAll(ctx context.Context, subscriber string) error {
	if !w.IsRunning() {
		return errNotRunning
	}

	if err := w.ws.UnsubscribeAll(ctx); err != nil {
		return err
	}

	w.mtx.Lock()
	w.subscriptions = make(map[string]chan ctypes.ResultEvent)
	w.mtx.Unlock()

	return nil
}

// After being reconnected, it is necessary to redo subscription to server
// otherwise no data will be automatically received.
func (w *WSEvents) redoSubscriptionsAfter(d time.Duration) {
	time.Sleep(d)

	w.mtx.RLock()
	defer w.mtx.RUnlock()
	for q := range w.subscriptions {
		err := w.ws.Subscribe(context.Background(), q)
		if err != nil {
			w.Logger.Error("Failed to resubscribe", "err", err)
		}
	}
}

func isErrAlreadySubscribed(err error) bool {
	return strings.Contains(err.Error(), tmpubsub.ErrAlreadySubscribed.Error())
}

func (w *WSEvents) eventListener() {
	for {
		select {
		case resp, ok := <-w.ws.ResponsesCh:
			if !ok {
				return
			}

			if resp.Error != nil {
				w.Logger.Error("WS error", "err", resp.Error.Error())
				// Error can be ErrAlreadySubscribed or max client (subscriptions per
				// client) reached or Tendermint exited.
				// We can ignore ErrAlreadySubscribed, but need to retry in other
				// cases.
				if !isErrAlreadySubscribed(resp.Error) {
					// Resubscribe after 1 second to give Tendermint time to restart (if
					// crashed).
					w.redoSubscriptionsAfter(1 * time.Second)
				}
				continue
			}

			result := new(ctypes.ResultEvent)
			err := tmjson.Unmarshal(resp.Result, result)
			if err != nil {
				w.Logger.Error("failed to unmarshal response", "err", err)
				continue
			}

			w.mtx.RLock()
			if out, ok := w.subscriptions[result.Query]; ok {
				if cap(out) == 0 {
					out <- *result
				} else {
					select {
					case out <- *result:
					default:
						w.Logger.Error("wanted to publish ResultEvent, but out channel is full", "result", result, "query", result.Query)
					}
				}
			}
			w.mtx.RUnlock()
		case <-w.Quit():
			return
		}
	}
}

// WSClient is a JSON-RPC client, which uses WebSocket for communication with
// the remote server.
//
// WSClient is safe for concurrent use by multiple goroutines.
type WSClient struct { // nolint: maligned
	conn *websocket.Conn

	Address  string // IP:PORT or /path/to/socket
	Endpoint string // /websocket/url/endpoint
	Header   http.Header
	Dialer   func(context.Context, string, string) (net.Conn, error)

	// Single user facing channel to read RPCResponses from, closed only when the
	// client is being stopped.
	ResponsesCh chan types.RPCResponse

	// Callback, which will be called each time after successful reconnect.
	onReconnect func()

	// internal channels
	send            chan types.RPCRequest // user requests
	backlog         chan types.RPCRequest // stores a single user request received during a conn failure
	reconnectAfter  chan error            // reconnect requests
	readRoutineQuit chan struct{}         // a way for readRoutine to close writeRoutine

	// Maximum reconnect attempts (0 or greater; default: 25).
	maxReconnectAttempts int

	// Support both ws and wss protocols
	protocol string

	wg sync.WaitGroup

	mtx            tmsync.RWMutex
	sentLastPingAt time.Time
	reconnecting   bool
	nextReqID      int
	// sentIDs        map[types.JSONRPCIntID]bool // IDs of the requests currently in flight

	// Time allowed to write a message to the server. 0 means block until operation succeeds.
	writeWait time.Duration

	// Time allowed to read the next message from the server. 0 means block until operation succeeds.
	readWait time.Duration

	// Send pings to server with this period. Must be less than readWait. If 0, no pings will be sent.
	pingPeriod time.Duration

	service.BaseService

	// Time between sending a ping and receiving a pong. See
	// https://godoc.org/github.com/rcrowley/go-metrics#Timer.
	PingPongLatencyTimer metrics.Timer
}

// NewWS returns a new client. See the commentary on the func(*WSClient)
// functions for a detailed description of how to configure ping period and
// pong wait time. The endpoint argument must begin with a `/`.
// An error is returned on invalid remote. The function panics when remote is nil.
func NewWS(remoteAddr, endpoint string, header http.Header, options ...func(*WSClient)) (*WSClient, error) {
	parsedURL, err := newParsedURL(remoteAddr)
	if err != nil {
		return nil, err
	}
	// default to ws protocol, unless wss is explicitly specified
	if parsedURL.Scheme != protoWSS {
		parsedURL.Scheme = protoWS
	}

	dialFn, err := makeHTTPDialer(remoteAddr)
	if err != nil {
		return nil, err
	}

	c := &WSClient{
		Address:              parsedURL.Host + parsedURL.Path,
		Dialer:               dialFn,
		Endpoint:             endpoint,
		Header:               header,
		PingPongLatencyTimer: metrics.NewTimer(),

		maxReconnectAttempts: defaultMaxReconnectAttempts,
		readWait:             defaultReadWait,
		writeWait:            defaultWriteWait,
		pingPeriod:           defaultPingPeriod,
		protocol:             parsedURL.Scheme,

		// sentIDs: make(map[types.JSONRPCIntID]bool),
	}
	c.BaseService = *service.NewBaseService(nil, "WSClient", c)
	for _, option := range options {
		option(c)
	}
	return c, nil
}

// MaxReconnectAttempts sets the maximum number of reconnect attempts before returning an error.
// It should only be used in the constructor and is not Goroutine-safe.
func MaxReconnectAttempts(max int) func(*WSClient) {
	return func(c *WSClient) {
		c.maxReconnectAttempts = max
	}
}

// ReadWait sets the amount of time to wait before a websocket read times out.
// It should only be used in the constructor and is not Goroutine-safe.
func ReadWait(readWait time.Duration) func(*WSClient) {
	return func(c *WSClient) {
		c.readWait = readWait
	}
}

// WriteWait sets the amount of time to wait before a websocket write times out.
// It should only be used in the constructor and is not Goroutine-safe.
func WriteWait(writeWait time.Duration) func(*WSClient) {
	return func(c *WSClient) {
		c.writeWait = writeWait
	}
}

// PingPeriod sets the duration for sending websocket pings.
// It should only be used in the constructor - not Goroutine-safe.
func PingPeriod(pingPeriod time.Duration) func(*WSClient) {
	return func(c *WSClient) {
		c.pingPeriod = pingPeriod
	}
}

// OnReconnect sets the callback, which will be called every time after
// successful reconnect.
func OnReconnect(cb func()) func(*WSClient) {
	return func(c *WSClient) {
		c.onReconnect = cb
	}
}

// String returns WS client full address.
func (c *WSClient) String() string {
	return fmt.Sprintf("WSClient{%s (%s)}", c.Address, c.Endpoint)
}

// OnStart implements service.Service by dialing a server and creating read and
// write routines.
func (c *WSClient) OnStart() error {
	err := c.dial(c.Header)
	if err != nil {
		return err
	}

	c.ResponsesCh = make(chan types.RPCResponse)

	c.send = make(chan types.RPCRequest)
	// 1 additional error may come from the read/write
	// goroutine depending on which failed first.
	c.reconnectAfter = make(chan error, 1)
	// capacity for 1 request. a user won't be able to send more because the send
	// channel is unbuffered.
	c.backlog = make(chan types.RPCRequest, 1)

	c.startReadWriteRoutines()
	go c.reconnectRoutine()

	return nil
}

// Stop overrides service.Service#Stop. There is no other way to wait until Quit
// channel is closed.
func (c *WSClient) Stop() error {
	if err := c.BaseService.Stop(); err != nil {
		return err
	}
	// only close user-facing channels when we can't write to them
	c.wg.Wait()
	close(c.ResponsesCh)

	return nil
}

// IsReconnecting returns true if the client is reconnecting right now.
func (c *WSClient) IsReconnecting() bool {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.reconnecting
}

// IsActive returns true if the client is running and not reconnecting.
func (c *WSClient) IsActive() bool {
	return c.IsRunning() && !c.IsReconnecting()
}

// Send the given RPC request to the server. Results will be available on
// ResponsesCh, errors, if any, on ErrorsCh. Will block until send succeeds or
// ctx.Done is closed.
func (c *WSClient) Send(ctx context.Context, request types.RPCRequest) error {
	select {
	case c.send <- request:
		c.Logger.Info("sent a request", "req", request)
		//c.mtx.Lock()
		//c.sentIDs[request.ID.(types.JSONRPCIntID)] = true
		//c.mtx.Unlock()
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Call enqueues a call request onto the Send queue. Requests are JSON encoded.
func (c *WSClient) Call(ctx context.Context, method string, params map[string]interface{}) error {
	request, err := types.MapToRequest(c.nextRequestID(), method, params)
	if err != nil {
		return err
	}
	return c.Send(ctx, request)
}

// CallWithArrayParams enqueues a call request onto the Send queue. Params are
// in a form of array (e.g. []interface{}{"abcd"}). Requests are JSON encoded.
func (c *WSClient) CallWithArrayParams(ctx context.Context, method string, params []interface{}) error {
	request, err := types.ArrayToRequest(c.nextRequestID(), method, params)
	if err != nil {
		return err
	}
	return c.Send(ctx, request)
}

// Private methods

func (c *WSClient) nextRequestID() types.JSONRPCIntID {
	c.mtx.Lock()
	id := c.nextReqID
	c.nextReqID++
	c.mtx.Unlock()
	return types.JSONRPCIntID(id)
}

func (c *WSClient) dial(header http.Header) error {
	dialer := &websocket.Dialer{
		NetDialContext: c.Dialer,
		Proxy:          http.ProxyFromEnvironment,
	}
	conn, _, err := dialer.Dial(c.protocol+"://"+c.Address+c.Endpoint, header) // nolint:bodyclose
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

// reconnect tries to redial up to maxReconnectAttempts with exponential
// backoff.
func (c *WSClient) reconnect() error {
	attempt := 0

	c.mtx.Lock()
	c.reconnecting = true
	c.mtx.Unlock()
	defer func() {
		c.mtx.Lock()
		c.reconnecting = false
		c.mtx.Unlock()
	}()

	for {
		jitter := time.Duration(tmrand.Float64() * float64(time.Second)) // 1s == (1e9 ns)
		backoffDuration := jitter + ((1 << uint(attempt)) * time.Second)

		c.Logger.Info("reconnecting", "attempt", attempt+1, "backoff_duration", backoffDuration)
		time.Sleep(backoffDuration)

		err := c.dial(c.Header)
		if err != nil {
			c.Logger.Error("failed to redial", "err", err)
		} else {
			c.Logger.Info("reconnected")
			if c.onReconnect != nil {
				go c.onReconnect()
			}
			return nil
		}

		attempt++

		if attempt > c.maxReconnectAttempts {
			return fmt.Errorf("reached maximum reconnect attempts: %w", err)
		}
	}
}

func (c *WSClient) startReadWriteRoutines() {
	c.wg.Add(2)
	c.readRoutineQuit = make(chan struct{})
	go c.readRoutine()
	go c.writeRoutine()
}

func (c *WSClient) processBacklog() error {
	select {
	case request := <-c.backlog:
		if c.writeWait > 0 {
			if err := c.conn.SetWriteDeadline(time.Now().Add(c.writeWait)); err != nil {
				c.Logger.Error("failed to set write deadline", "err", err)
			}
		}
		if err := c.conn.WriteJSON(request); err != nil {
			c.Logger.Error("failed to resend request", "err", err)
			c.reconnectAfter <- err
			// requeue request
			c.backlog <- request
			return err
		}
		c.Logger.Info("resend a request", "req", request)
	default:
	}
	return nil
}

func (c *WSClient) reconnectRoutine() {
	for {
		select {
		case originalError := <-c.reconnectAfter:
			// wait until writeRoutine and readRoutine finish
			c.wg.Wait()
			if err := c.reconnect(); err != nil {
				c.Logger.Error("failed to reconnect", "err", err, "original_err", originalError)
				if err = c.Stop(); err != nil {
					c.Logger.Error("failed to stop conn", "error", err)
				}

				return
			}
			// drain reconnectAfter
		LOOP:
			for {
				select {
				case <-c.reconnectAfter:
				default:
					break LOOP
				}
			}
			err := c.processBacklog()
			if err == nil {
				c.startReadWriteRoutines()
			}

		case <-c.Quit():
			return
		}
	}
}

// The client ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *WSClient) writeRoutine() {
	var ticker *time.Ticker
	if c.pingPeriod > 0 {
		// ticker with a predefined period
		ticker = time.NewTicker(c.pingPeriod)
	} else {
		// ticker that never fires
		ticker = &time.Ticker{C: make(<-chan time.Time)}
	}

	defer func() {
		ticker.Stop()
		c.conn.Close()
		// err != nil {
		// ignore error; it will trigger in tests
		// likely because it's closing an already closed connection
		// }
		c.wg.Done()
	}()

	for {
		select {
		case request := <-c.send:
			if c.writeWait > 0 {
				if err := c.conn.SetWriteDeadline(time.Now().Add(c.writeWait)); err != nil {
					c.Logger.Error("failed to set write deadline", "err", err)
				}
			}
			if err := c.conn.WriteJSON(request); err != nil {
				c.Logger.Error("failed to send request", "err", err)
				c.reconnectAfter <- err
				// add request to the backlog, so we don't lose it
				c.backlog <- request
				return
			}
		case <-ticker.C:
			if c.writeWait > 0 {
				if err := c.conn.SetWriteDeadline(time.Now().Add(c.writeWait)); err != nil {
					c.Logger.Error("failed to set write deadline", "err", err)
				}
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				c.Logger.Error("failed to write ping", "err", err)
				c.reconnectAfter <- err
				return
			}
			c.mtx.Lock()
			c.sentLastPingAt = time.Now()
			c.mtx.Unlock()
			c.Logger.Debug("sent ping")
		case <-c.readRoutineQuit:
			return
		case <-c.Quit():
			if err := c.conn.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			); err != nil {
				c.Logger.Error("failed to write message", "err", err)
			}
			return
		}
	}
}

// The client ensures that there is at most one reader to a connection by
// executing all reads from this goroutine.
func (c *WSClient) readRoutine() {
	defer func() {
		c.conn.Close()
		// err != nil {
		// ignore error; it will trigger in tests
		// likely because it's closing an already closed connection
		// }
		c.wg.Done()
	}()

	c.conn.SetPongHandler(func(string) error {
		// gather latency stats
		c.mtx.RLock()
		t := c.sentLastPingAt
		c.mtx.RUnlock()
		c.PingPongLatencyTimer.UpdateSince(t)

		c.Logger.Debug("got pong")
		return nil
	})

	for {
		// reset deadline for every message type (control or data)
		if c.readWait > 0 {
			if err := c.conn.SetReadDeadline(time.Now().Add(c.readWait)); err != nil {
				c.Logger.Error("failed to set read deadline", "err", err)
			}
		}
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if !websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
				return
			}

			c.Logger.Error("failed to read response", "err", err)
			close(c.readRoutineQuit)
			c.reconnectAfter <- err
			return
		}

		var response types.RPCResponse
		err = json.Unmarshal(data, &response)
		if err != nil {
			c.Logger.Error("failed to parse response", "err", err, "data", string(data))
			continue
		}

		if err = validateResponseID(response.ID); err != nil {
			c.Logger.Error("error in response ID", "id", response.ID, "err", err)
			continue
		}

		// TODO: events resulting from /subscribe do not work with ->
		// because they are implemented as responses with the subscribe request's
		// ID. According to the spec, they should be notifications (requests
		// without IDs).
		// https://github.com/tendermint/tendermint/issues/2949
		// c.mtx.Lock()
		// if _, ok := c.sentIDs[response.ID.(types.JSONRPCIntID)]; !ok {
		// 	c.Logger.Error("unsolicited response ID", "id", response.ID, "expected", c.sentIDs)
		// 	c.mtx.Unlock()
		// 	continue
		// }
		// delete(c.sentIDs, response.ID.(types.JSONRPCIntID))
		// c.mtx.Unlock()
		// Combine a non-blocking read on BaseService.Quit with a non-blocking write on ResponsesCh to avoid blocking
		// c.wg.Wait() in c.Stop(). Note we rely on Quit being closed so that it sends unlimited Quit signals to stop
		// both readRoutine and writeRoutine

		c.Logger.Info("got response", "id", response.ID, "result", fmt.Sprintf("%X", response.Result))

		select {
		case <-c.Quit():
		case c.ResponsesCh <- response:
		}
	}
}

// Predefined methods

// Subscribe to a query. Note the server must have a "subscribe" route
// defined.
func (c *WSClient) Subscribe(ctx context.Context, query string) error {
	params := map[string]interface{}{"query": query}
	return c.Call(ctx, "subscribe", params)
}

// Unsubscribe from a query. Note the server must have a "unsubscribe" route
// defined.
func (c *WSClient) Unsubscribe(ctx context.Context, query string) error {
	params := map[string]interface{}{"query": query}
	return c.Call(ctx, "unsubscribe", params)
}

// UnsubscribeAll from all. Note the server must have a "unsubscribe_all" route
// defined.
func (c *WSClient) UnsubscribeAll(ctx context.Context) error {
	params := map[string]interface{}{}
	return c.Call(ctx, "unsubscribe_all", params)
}

// Parsed URL structure
type parsedURL struct {
	url.URL
}

// Parse URL and set defaults
func newParsedURL(remoteAddr string) (*parsedURL, error) {
	u, err := url.Parse(remoteAddr)
	if err != nil {
		return nil, err
	}

	// default to tcp if nothing specified
	if u.Scheme == "" {
		u.Scheme = protoTCP
	}

	return &parsedURL{*u}, nil
}

// Change protocol to HTTP for unknown protocols and TCP protocol - useful for RPC connections
func (u *parsedURL) SetDefaultSchemeHTTP() {
	// protocol to use for http operations, to support both http and https
	switch u.Scheme {
	case protoHTTP, protoHTTPS, protoWS, protoWSS:
		// known protocols not changed
	default:
		// default to http for unknown protocols (ex. tcp)
		u.Scheme = protoHTTP
	}
}

// Get full address without the protocol - useful for Dialer connections
func (u parsedURL) GetHostWithPath() string {
	// Remove protocol, userinfo and # fragment, assume opaque is empty
	return u.Host + u.EscapedPath()
}

// Get a trimmed address with protocol - useful as address in RPC connections
func (u parsedURL) GetTrimmedURL() string {
	return u.Scheme + "://" + u.GetTrimmedHostWithPath()
}

// Get a address with protocol - useful as address in RPC connections
func (u parsedURL) GetURL() string {
	return u.Scheme + "://" + u.GetHostWithPath()
}

// Get a trimmed address - useful for WS connections
func (u parsedURL) GetTrimmedHostWithPath() string {
	// replace / with . for http requests (kvstore domain)
	return strings.ReplaceAll(u.GetHostWithPath(), "/", ".")
}
