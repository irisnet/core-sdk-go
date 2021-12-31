package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmjson "github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/libs/service"
	"github.com/tendermint/tendermint/rpc/client"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/rpc/jsonrpc/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

var errNotRunning = errors.New("client is not running. Use .Start() method to start")

var _ service.Service = (*JSONRpcClient)(nil)

type JSONRpcClient struct {
	remote string
	client *http.Client
	header http.Header
	*WSEvents
}

func NewJSONRpcClient(rpcAddr, wsAddr string, endpoint string, timeout uint, header http.Header) (JSONRpcClient, error) {
	if wsAddr == "" {
		wsAddr = rpcAddr
	}
	wsEvents, err := newWSEvents(wsAddr, endpoint, header)
	if err != nil {
		return JSONRpcClient{}, err
	}

	httpClient, err := DefaultHTTPClient(rpcAddr)
	if err != nil {
		return JSONRpcClient{}, err
	}

	parsedURL, err := newParsedURL(rpcAddr)
	if err != nil {
		return JSONRpcClient{}, fmt.Errorf("invalid rpcAddr %s: %s", rpcAddr, err)
	}
	parsedURL.SetDefaultSchemeHTTP()

	httpClient.Timeout = time.Duration(timeout) * time.Second
	return JSONRpcClient{
		remote:   parsedURL.GetTrimmedURL(),
		client:   httpClient,
		header:   header,
		WSEvents: wsEvents,
	}, nil
}

func (c JSONRpcClient) Quit() <-chan struct{} {
	panic("implement me")
}

func (c JSONRpcClient) String() string {
	panic("implement me")
}

func (c JSONRpcClient) SetLogger(logger log.Logger) {
	panic("implement me")
}

func (c JSONRpcClient) ABCIInfo(ctx context.Context) (*ctypes.ResultABCIInfo, error) {
	panic("implement me")
}

func (c JSONRpcClient) ABCIQuery(ctx context.Context, path string, data tmbytes.HexBytes) (*ctypes.ResultABCIQuery, error) {
	return c.ABCIQueryWithOptions(ctx, path, data, rpcclient.DefaultABCIQueryOptions)
}

func (c JSONRpcClient) ABCIQueryWithOptions(ctx context.Context, path string, data tmbytes.HexBytes, opts client.ABCIQueryOptions) (*ctypes.ResultABCIQuery, error) {
	result := new(ctypes.ResultABCIQuery)
	_, err := c.Call(ctx, "abci_query", map[string]interface{}{"path": path, "data": data, "height": strconv.FormatInt(opts.Height, 10), "prove": opts.Prove}, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c JSONRpcClient) BroadcastTxCommit(ctx context.Context, tx tmtypes.Tx) (*ctypes.ResultBroadcastTxCommit, error) {
	result := new(ctypes.ResultBroadcastTxCommit)
	_, err := c.Call(ctx, "broadcast_tx_commit", map[string]interface{}{"tx": tx}, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c JSONRpcClient) BroadcastTxAsync(ctx context.Context, tx tmtypes.Tx) (*ctypes.ResultBroadcastTx, error) {
	return c.broadcastTX(ctx, "broadcast_tx_async", tx)
}

func (c JSONRpcClient) BroadcastTxSync(ctx context.Context, tx tmtypes.Tx) (*ctypes.ResultBroadcastTx, error) {
	return c.broadcastTX(ctx, "broadcast_tx_sync", tx)
}

func (c JSONRpcClient) Genesis(ctx context.Context) (*ctypes.ResultGenesis, error) {
	result := new(ctypes.ResultGenesis)
	_, err := c.Call(ctx, "genesis", map[string]interface{}{}, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c JSONRpcClient) BlockchainInfo(ctx context.Context, minHeight, maxHeight int64) (*ctypes.ResultBlockchainInfo, error) {
	result := new(ctypes.ResultBlockchainInfo)
	_, err := c.Call(ctx, "blockchain",
		map[string]interface{}{"minHeight": minHeight, "maxHeight": maxHeight},
		result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c JSONRpcClient) NetInfo(ctx context.Context) (*ctypes.ResultNetInfo, error) {
	result := new(ctypes.ResultNetInfo)
	_, err := c.Call(ctx, "net_info", map[string]interface{}{}, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c JSONRpcClient) DumpConsensusState(ctx context.Context) (*ctypes.ResultDumpConsensusState, error) {
	result := new(ctypes.ResultDumpConsensusState)
	_, err := c.Call(ctx, "dump_consensus_state", map[string]interface{}{}, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c JSONRpcClient) ConsensusState(ctx context.Context) (*ctypes.ResultConsensusState, error) {
	result := new(ctypes.ResultConsensusState)
	_, err := c.Call(ctx, "consensus_state", map[string]interface{}{}, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c JSONRpcClient) ConsensusParams(ctx context.Context, height *int64) (*ctypes.ResultConsensusParams, error) {
	result := new(ctypes.ResultConsensusParams)
	params := make(map[string]interface{})
	if height != nil {
		params["height"] = height
	}
	_, err := c.Call(ctx, "consensus_params", params, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c JSONRpcClient) Health(ctx context.Context) (*ctypes.ResultHealth, error) {
	result := new(ctypes.ResultHealth)
	_, err := c.Call(ctx, "health", map[string]interface{}{}, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c JSONRpcClient) Block(ctx context.Context, height *int64) (*ctypes.ResultBlock, error) {
	result := new(ctypes.ResultBlock)
	params := make(map[string]interface{})
	if height != nil {
		params["height"] = strconv.Itoa(int(*height))
	}
	_, err := c.Call(ctx, "block", params, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c JSONRpcClient) BlockByHash(ctx context.Context, hash []byte) (*ctypes.ResultBlock, error) {
	result := new(ctypes.ResultBlock)
	params := map[string]interface{}{
		"hash": hash,
	}
	_, err := c.Call(ctx, "block_by_hash", params, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c JSONRpcClient) BlockResults(ctx context.Context, height *int64) (*ctypes.ResultBlockResults, error) {
	result := new(ctypes.ResultBlockResults)
	params := make(map[string]interface{})
	if height != nil {
		params["height"] = strconv.Itoa(int(*height))
	}
	_, err := c.Call(ctx, "block_results", params, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c JSONRpcClient) Commit(ctx context.Context, height *int64) (*ctypes.ResultCommit, error) {
	result := new(ctypes.ResultCommit)
	params := make(map[string]interface{})
	if height != nil {
		params["height"] = height
	}
	_, err := c.Call(ctx, "commit", params, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c JSONRpcClient) Validators(ctx context.Context, height *int64, page, perPage *int) (*ctypes.ResultValidators, error) {
	result := new(ctypes.ResultValidators)
	params := make(map[string]interface{})
	if page != nil {
		params["page"] = page
	}
	if perPage != nil {
		params["per_page"] = perPage
	}
	if height != nil {
		params["height"] = height
	}
	_, err := c.Call(ctx, "validators", params, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c JSONRpcClient) Tx(ctx context.Context, hash []byte, prove bool) (*ctypes.ResultTx, error) {
	result := new(ctypes.ResultTx)
	params := map[string]interface{}{
		"hash":  hash,
		"prove": prove,
	}
	_, err := c.Call(ctx, "tx", params, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c JSONRpcClient) TxSearch(ctx context.Context, query string, prove bool, page, perPage *int, orderBy string) (*ctypes.ResultTxSearch, error) {
	result := new(ctypes.ResultTxSearch)
	params := map[string]interface{}{
		"query":    query,
		"prove":    prove,
		"order_by": orderBy,
	}
	if page != nil {
		params["page"] = strconv.Itoa(*page)
	}
	if perPage != nil {
		params["per_page"] = strconv.Itoa(*perPage)
	}
	_, err := c.Call(ctx, "tx_search", params, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c JSONRpcClient) Status(ctx context.Context) (*ctypes.ResultStatus, error) {
	result := new(ctypes.ResultStatus)
	_, err := c.Call(ctx, "status", map[string]interface{}{}, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c JSONRpcClient) BroadcastEvidence(ctx context.Context, ev tmtypes.Evidence) (*ctypes.ResultBroadcastEvidence, error) {
	result := new(ctypes.ResultBroadcastEvidence)
	_, err := c.Call(ctx, "broadcast_evidence", map[string]interface{}{"evidence": ev}, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c JSONRpcClient) UnconfirmedTxs(ctx context.Context, limit *int) (*ctypes.ResultUnconfirmedTxs, error) {
	result := new(ctypes.ResultUnconfirmedTxs)
	params := make(map[string]interface{})
	if limit != nil {
		params["limit"] = limit
	}
	_, err := c.Call(ctx, "unconfirmed_txs", params, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c JSONRpcClient) NumUnconfirmedTxs(ctx context.Context) (*ctypes.ResultUnconfirmedTxs, error) {
	result := new(ctypes.ResultUnconfirmedTxs)
	_, err := c.Call(ctx, "num_unconfirmed_txs", map[string]interface{}{}, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c JSONRpcClient) CheckTx(ctx context.Context, tx tmtypes.Tx) (*ctypes.ResultCheckTx, error) {
	result := new(ctypes.ResultCheckTx)
	_, err := c.Call(ctx, "check_tx", map[string]interface{}{"tx": tx}, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *JSONRpcClient) mapToRequest(method string, params map[string]interface{}) ([]byte, error) {
	var paramsMap = make(map[string]interface{})
	paramsMap["jsonrpc"] = "2.0"
	paramsMap["id"] = 0
	paramsMap["method"] = method
	paramsMap["params"] = params
	return json.Marshal(paramsMap)
}
func (c *JSONRpcClient) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

func (c *JSONRpcClient) Call(ctx context.Context, method string, params map[string]interface{}, result interface{}) (interface{}, error) {
	requestBytes, err := c.mapToRequest(method, params)
	if err != nil {
		return nil, fmt.Errorf("request failed: %s", err.Error())
	}

	req, err := http.NewRequest(http.MethodPost, c.remote, bytes.NewReader(requestBytes))
	if err != nil {
		return nil, fmt.Errorf("request failed: %s", err.Error())
	}

	if c.header != nil {
		for h := range c.header {
			req.Header.Add(h, c.header.Get(h))
		}
	}

	httpResponse, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("post failed: %s", err.Error())
	}
	defer httpResponse.Body.Close()

	httpResponseBytes, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %s", err.Error())
	}
	rpcResponse := &types.RPCResponse{}
	if err = json.Unmarshal(httpResponseBytes, rpcResponse); err != nil {
		return nil, fmt.Errorf("error unmarshalling: %s", err.Error())
	}
	if rpcResponse.Error != nil {
		return nil, fmt.Errorf("request failed, code: %d, message: %s, data: %s", rpcResponse.Error.Code, rpcResponse.Error.Message, rpcResponse.Error.Data)
	}
	if err = tmjson.Unmarshal(rpcResponse.Result, result); err != nil {
		return nil, fmt.Errorf("error unmarshalling result: %s", err.Error())
	}
	return result, nil
}

func (c JSONRpcClient) broadcastTX(ctx context.Context, route string, tx tmtypes.Tx) (*ctypes.ResultBroadcastTx, error) {
	result := new(ctypes.ResultBroadcastTx)
	_, err := c.Call(ctx, route, map[string]interface{}{"tx": tx}, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

//-------------------------------------------------------------

func makeHTTPDialer(remoteAddr string) (func(string, string) (net.Conn, error), error) {
	u, err := newParsedURL(remoteAddr)
	if err != nil {
		return nil, err
	}

	protocol := u.Scheme

	// accept http(s) as an alias for tcp
	switch protocol {
	case protoHTTP, protoHTTPS:
		protocol = protoTCP
	}

	dialFn := func(proto, addr string) (net.Conn, error) {
		return net.Dial(protocol, u.GetHostWithPath())
	}

	return dialFn, nil
}

// DefaultHTTPClient is used to create an http client with some default parameters.
// We overwrite the http.Client.Dial so we can do http over tcp or unix.
// remoteAddr should be fully featured (eg. with tcp:// or unix://).
// An error will be returned in case of invalid remoteAddr.
func DefaultHTTPClient(remoteAddr string) (*http.Client, error) {
	dialFn, err := makeHTTPDialer(remoteAddr)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			// Set to true to prevent GZIP-bomb DoS attacks
			DisableCompression: true,
			Dial:               dialFn,
		},
	}

	return client, nil
}
