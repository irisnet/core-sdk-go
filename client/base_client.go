package client

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/avast/retry-go"
	grpc1 "github.com/gogo/protobuf/grpc"
	"github.com/gogo/protobuf/proto"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/libs/log"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/irisnet/core-sdk-go/cache"
	"github.com/irisnet/core-sdk-go/codec"
	sdklog "github.com/irisnet/core-sdk-go/log"
	"github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/core-sdk-go/types/errors"
	"github.com/irisnet/core-sdk-go/types/tx"
)

const (
	concurrency       = 16
	cacheCapacity     = 100
	cacheExpirePeriod = 1 * time.Minute
	tryThreshold      = 3
	maxBatch          = 100
)

type baseClient struct {
	types.TmClient
	types.TokenManager
	types.KeyManager
	cfg            *types.ClientConfig
	encodingConfig types.EncodingConfig
	l              *locker
	AccountQuery
}

// NewBaseClient return the baseClient for every sub modules
func NewBaseClient(cfg types.ClientConfig, encodingConfig types.EncodingConfig, logger log.Logger) types.BaseClient {
	if logger == nil {
		logger = sdklog.NewLogger(sdklog.Config{
			Format: sdklog.FormatText,
			Level:  cfg.Level,
		})
	}

	base := baseClient{
		TmClient:       NewRPCClient(cfg.NodeURI, encodingConfig.TxConfig.TxDecoder(), logger, cfg.Timeout),
		cfg:            &cfg,
		encodingConfig: encodingConfig,
		l:              NewLocker(concurrency).setLogger(logger),
		TokenManager:   cfg.TokenManager,
	}
	base.KeyManager = NewKeyManager(cfg.KeyDAO, cfg.Algo)

	c := cache.NewCache(cacheCapacity, cfg.Cached)
	base.AccountQuery = AccountQuery{
		Queries:    base,
		GRPCClient: NewGRPCClient(cfg.GRPCAddr),
		Logger:     logger,
		Cache:      c,
		cdc:        encodingConfig.Codec,
		Km:         base.KeyManager,
		expiration: cacheExpirePeriod,
	}
	return &base
}

func (base *baseClient) Logger() log.Logger {
	return base.AccountQuery.Logger
}

func (base *baseClient) SetLogger(logger log.Logger) {
	base.AccountQuery.Logger = logger
}

// Codec returns codec.
func (base *baseClient) Marshaler() codec.Codec {
	return base.encodingConfig.Codec
}

func (base *baseClient) GenConn() (grpc1.ClientConn, error) {
	return base.AccountQuery.GenConn()
}

func (base *baseClient) BuildTxHash(msg []types.Msg, baseTx types.BaseTx) (string, error) {
	txByte, _, err := base.buildTx(msg, baseTx)
	if err != nil {
		return "", errors.Wrap(errors.ErrTodo, err.Error())
	}
	return strings.ToUpper(hex.EncodeToString(tmhash.Sum(txByte))), nil
}

func (base *baseClient) BuildAndSign(msg []types.Msg, baseTx types.BaseTx) ([]byte, error) {
	builder, err := base.prepare(baseTx)
	if err != nil {
		return nil, errors.Wrap(errors.ErrTodo, err.Error())
	}

	txByte, err := builder.BuildAndSign(baseTx.From, msg, true)
	if err != nil {
		return nil, errors.Wrap(errors.ErrTodo, err.Error())
	}

	base.Logger().Debug("sign transaction success")
	return txByte, nil
}

func (base *baseClient) BuildAndSendWithAccount(addr string, accountNumber, sequence uint64, msg []types.Msg, baseTx types.BaseTx) (ctypes.ResultTx, error) {
	txByte, ctx, err := base.buildTxWithAccount(addr, accountNumber, sequence, msg, baseTx)
	if err != nil {
		return ctypes.ResultTx{}, err
	}

	valid, err := base.ValidateTxSize(len(txByte), msg)
	if err != nil {
		return ctypes.ResultTx{}, err
	}
	if !valid {
		base.Logger().Debug("tx is too large")
		// filter out transactions that have been sent
		// reset the maximum number of msg in each transaction
		//batch = batch / 2
		return ctypes.ResultTx{}, errors.ErrTxTooLarge
	}
	return base.broadcastTx(txByte, ctx.Mode())
}

func (base *baseClient) BuildAndSend(msg []types.Msg, baseTx types.BaseTx) (ctypes.ResultTx, error) {
	var res ctypes.ResultTx
	var address string

	// lock the account
	base.l.Lock(baseTx.From)
	defer base.l.Unlock(baseTx.From)

	retryableFunc := func() error {
		txByte, ctx, e := base.buildTx(msg, baseTx)
		if e != nil {
			return e
		}

		if res, e = base.broadcastTx(txByte, ctx.Mode()); e != nil {
			address = ctx.Address()
			return e
		}
		return nil
	}

	retryIfFunc := func(err error) bool {
		return errors.Code(err) == errors.ErrWrongSequence.Code()
	}

	onRetryFunc := func(n uint, err error) {
		_ = base.removeCache(address)
		base.Logger().Error(
			"wrong sequence, will retry",
			"address", address,
			"attempts", n,
			"err", err.Error(),
		)
	}

	if err := retry.Do(
		retryableFunc,
		retry.Attempts(tryThreshold),
		retry.RetryIf(retryIfFunc),
		retry.OnRetry(onRetryFunc),
	); err != nil {
		return res, errors.Wrap(errors.ErrTodo, err.Error())
	}

	return res, nil
}

func (base *baseClient) SendBatch(msgs types.Msgs, baseTx types.BaseTx) (rs []ctypes.ResultTx, err error) {
	if msgs == nil || len(msgs) == 0 {
		return rs, errors.Wrapf(errors.ErrTodo, "must have at least one message in list")
	}

	defer errors.CatchPanic(func(errMsg string) {
		base.Logger().Error("broadcast msg failed", "errMsg", errMsg)
	})
	// validate msg
	for _, m := range msgs {
		if err := m.ValidateBasic(); err != nil {
			return rs, errors.Wrap(errors.ErrTodo, err.Error())
		}
	}
	base.Logger().Debug("validate msg success")

	// lock the account
	base.l.Lock(baseTx.From)
	defer base.l.Unlock(baseTx.From)

	var address string
	var batch = maxBatch

	retryableFunc := func() error {
		for i, ms := range types.SubArray(batch, msgs) {
			mss := ms.(types.Msgs)
			txByte, ctx, err := base.buildTx(mss, baseTx)
			if err != nil {
				return err
			}

			valid, err := base.ValidateTxSize(len(txByte), mss)
			if err != nil {
				return err
			}
			if !valid {
				base.Logger().Debug("tx is too large", "msgsLength", batch)
				// filter out transactions that have been sent
				msgs = msgs[i*batch:]
				// reset the maximum number of msg in each transaction
				batch = batch / 2
				return errors.ErrTxTooLarge
			}
			res, err := base.broadcastTx(txByte, ctx.Mode())
			if err != nil {
				address = ctx.Address()
				return err
			}
			rs = append(rs, res)
		}
		return nil
	}

	retryIf := func(err error) bool {
		return errors.Code(err) == errors.ErrInvalidSequence.Code() || errors.Code(err) == errors.ErrTxTooLarge.Code()
	}

	onRetry := func(n uint, err error) {
		_ = base.removeCache(address)
		base.Logger().Error(
			"wrong sequence, will retry",
			"address", address,
			"attempts", n,
			"err", err.Error(),
		)
	}

	_ = retry.Do(retryableFunc,
		retry.Attempts(tryThreshold),
		retry.RetryIf(retryIf),
		retry.OnRetry(onRetry),
	)
	return rs, nil
}

func (base baseClient) QueryWithResponse(path string, data interface{}, result types.Response) error {
	res, err := base.Query(path, data)
	if err != nil {
		return err
	}

	if err := base.encodingConfig.Codec.UnmarshalJSON(res, result.(proto.Message)); err != nil {
		return err
	}

	return nil
}

func (base baseClient) Query(path string, data interface{}) ([]byte, error) {
	var bz []byte
	var err error
	if data != nil {
		bz, err = base.encodingConfig.Codec.MarshalJSON(data.(proto.Message))
		if err != nil {
			return nil, err
		}
	}

	opts := rpcclient.ABCIQueryOptions{
		// Height: cliCtx.Height,
		Prove: false,
	}
	result, err := base.ABCIQueryWithOptions(context.Background(), path, bz, opts)
	if err != nil {
		return nil, err
	}

	resp := result.Response
	if !resp.IsOK() {
		return nil, fmt.Errorf(resp.Log)
	}

	return resp.Value, nil
}

func (base baseClient) QueryStore(key types.HexBytes, storeName string, height int64, prove bool) (res abci.ResponseQuery, err error) {
	path := fmt.Sprintf("/store/%s/%s", storeName, "key")
	opts := rpcclient.ABCIQueryOptions{
		Prove:  prove,
		Height: height,
	}

	result, err := base.ABCIQueryWithOptions(context.Background(), path, key, opts)
	if err != nil {
		return res, err
	}

	resp := result.Response
	if !resp.IsOK() {
		return res, fmt.Errorf(resp.Log)
	}
	return resp, nil
}

func (base *baseClient) prepare(baseTx types.BaseTx) (*types.Factory, error) {
	factory := types.NewFactory().
		WithChainID(base.cfg.ChainID).
		WithKeyManager(base.AccountQuery.Km).
		WithMode(base.cfg.Mode).
		WithSimulateAndExecute(baseTx.SimulateAndExecute).
		WithGas(base.cfg.Gas).
		WithGasAdjustment(base.cfg.GasAdjustment).
		WithSignModeHandler(tx.MakeSignModeHandler(tx.DefaultSignModes)).
		WithTxConfig(base.encodingConfig.TxConfig).
		WithQueryFunc(base.QueryWithData)

	addr, err := base.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, err
	}
	factory.WithAddress(addr.String())

	account, err := base.QueryAndRefreshAccount(addr.String())
	if err != nil {
		return nil, err
	}
	factory.WithAccountNumber(account.AccountNumber).
		WithSequence(account.Sequence).
		WithPassword(baseTx.Password)

	if !baseTx.Fee.Empty() && baseTx.Fee.IsValid() {
		fees, err := base.TokenManager.ToMinCoin(baseTx.Fee...)
		if err != nil {
			return nil, err
		}
		factory.WithFee(fees)
	} else {
		fees, err := base.TokenManager.ToMinCoin(base.cfg.Fee...)
		if err != nil {
			panic(err)
		}
		factory.WithFee(fees)
	}

	if len(baseTx.Mode) > 0 {
		factory.WithMode(baseTx.Mode)
	}

	if baseTx.Gas > 0 {
		factory.WithGas(baseTx.Gas)
	}

	if baseTx.GasAdjustment > 0 {
		factory.WithGasAdjustment(baseTx.GasAdjustment)
	}

	if len(baseTx.Memo) > 0 {
		factory.WithMemo(baseTx.Memo)
	}
	return factory, nil
}

func (base *baseClient) prepareWithAccount(addr string, accountNumber, sequence uint64, baseTx types.BaseTx) (*types.Factory, error) {
	factory := types.NewFactory().
		WithChainID(base.cfg.ChainID).
		WithKeyManager(base.AccountQuery.Km).
		WithMode(base.cfg.Mode).
		WithSimulateAndExecute(baseTx.SimulateAndExecute).
		WithGas(base.cfg.Gas).
		WithSignModeHandler(tx.MakeSignModeHandler(tx.DefaultSignModes)).
		WithTxConfig(base.encodingConfig.TxConfig)

	factory.WithAddress(addr).
		WithAccountNumber(accountNumber).
		WithSequence(sequence).
		WithPassword(baseTx.Password)

	if !baseTx.Fee.Empty() && baseTx.Fee.IsValid() {
		fees, err := base.TokenManager.ToMinCoin(baseTx.Fee...)
		if err != nil {
			return nil, err
		}
		factory.WithFee(fees)
	} else {
		fees, err := base.TokenManager.ToMinCoin(base.cfg.Fee...)
		if err != nil {
			panic(err)
		}
		factory.WithFee(fees)
	}

	if len(baseTx.Mode) > 0 {
		factory.WithMode(baseTx.Mode)
	}

	if baseTx.Gas > 0 {
		factory.WithGas(baseTx.Gas)
	}

	if len(baseTx.Memo) > 0 {
		factory.WithMemo(baseTx.Memo)
	}
	return factory, nil
}

func (base *baseClient) ValidateTxSize(txSize int, msgs []types.Msg) (bool, error) {
	if uint64(txSize) > base.cfg.TxSizeLimit {
		return false, nil
	}
	return true, nil
}

type locker struct {
	shards []chan int
	size   int
	logger log.Logger
}

//NewLocker implement the function of lock, can lock resources according to conditions
func NewLocker(size int) *locker {
	shards := make([]chan int, size)
	for i := 0; i < size; i++ {
		shards[i] = make(chan int, 1)
	}
	return &locker{
		shards: shards,
		size:   size,
	}
}

func (l *locker) setLogger(logger log.Logger) *locker {
	l.logger = logger
	return l
}

func (l *locker) Lock(key string) {
	ch := l.getShard(key)
	ch <- 1
}

func (l *locker) Unlock(key string) {
	ch := l.getShard(key)
	<-ch
}

func (l *locker) getShard(key string) chan int {
	index := uint(l.indexFor(key)) % uint(l.size)
	return l.shards[index]
}

func (l *locker) indexFor(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}
