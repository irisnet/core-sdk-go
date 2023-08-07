// Package modules is to warpped the API provided by each module of IRITA
package client

import (
	"encoding/hex"
	"strings"
	"time"

	rpcclienthttp "github.com/tendermint/tendermint/rpc/client/http"

	"google.golang.org/grpc"

	"github.com/cosmos/cosmos-sdk/types"

	"github.com/avast/retry-go"
	commoncodec "github.com/cosmos/cosmos-sdk/codec"
	"github.com/irisnet/core-sdk-go/common/cache"
	commoncache "github.com/irisnet/core-sdk-go/common/cache"
	sdklog "github.com/irisnet/core-sdk-go/common/log"
	sdk "github.com/irisnet/core-sdk-go/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/libs/log"
	rpc "github.com/tendermint/tendermint/rpc/client"
)

const (
	concurrency       = 16
	cacheCapacity     = 100
	cacheExpirePeriod = 1 * time.Minute
	tryThreshold      = 3
	maxBatch          = 100
)

type baseClient struct {
	rpc.Client
	sdk.TokenManager
	cache.Cache
	log.Logger

	KeyManager sdk.KeyManager

	grpcConn   *grpc.ClientConn
	expiration time.Duration

	cfg            *sdk.ClientConfig
	encodingConfig sdk.EncodingConfig
	l              *locker
}

// NewBaseClient return the baseClient for every sub modules
func NewBaseClient(cfg sdk.ClientConfig, encodingConfig sdk.EncodingConfig) sdk.BaseClient {
	// create logger
	logger := sdklog.NewLogger(sdklog.Config{
		Format: sdklog.FormatText,
		Level:  cfg.Level,
	})

	grpcOptions := []grpc.DialOption{grpc.WithInsecure()}
	if len(cfg.GRPCOptions) > 0 {
		grpcOptions = append(grpcOptions, cfg.GRPCOptions...)
	}

	grpcConn, err := grpc.Dial(cfg.GRPCAddr, grpcOptions...)
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	cometbftClient, err := rpcclienthttp.NewWithTimeout(
		cfg.RPCAddr,
		cfg.WSAddr,
		cfg.Timeout)
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	base := baseClient{
		Client:         cometbftClient,
		grpcConn:       grpcConn,
		cfg:            &cfg,
		encodingConfig: encodingConfig,
		l:              NewLocker(concurrency).setLogger(logger),
		TokenManager:   cfg.TokenManager,
		KeyManager:     sdk.NewKeyManager(cfg.KeyDAO, cfg.Algo),
	}

	c := commoncache.NewCache(cacheCapacity, cfg.Cached)
	base.Cache = c

	base.Logger = logger

	return &base
}

func (base *baseClient) GrpcConn() *grpc.ClientConn {
	return base.grpcConn
}

func (base *baseClient) SetLogger(logger log.Logger) {
	base.Logger = logger
}

func (base *baseClient) Marshaler() commoncodec.Codec {
	return base.encodingConfig.Marshaler
}

func (base *baseClient) BuildTxHash(msg []types.Msg, baseTx sdk.BaseTx) (string, sdk.Error) {
	txByte, _, err := base.buildTx(msg, baseTx)
	if err != nil {
		return "", sdk.Wrap(err)
	}
	return strings.ToUpper(hex.EncodeToString(tmhash.Sum(txByte))), nil
}

func (base *baseClient) BuildAndSign(msg []types.Msg, baseTx sdk.BaseTx) ([]byte, sdk.Error) {
	builder, err := base.prepare(baseTx)
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	txByte, err := builder.BuildAndSign(baseTx.From, msg, false)
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	base.Debug("sign transaction success")
	return txByte, nil
}

func (base *baseClient) BuildAndSignWithAccount(addr string, accountNumber, sequence uint64, msg []types.Msg, baseTx sdk.BaseTx) ([]byte, sdk.Error) {
	txByte, _, err := base.buildTxWithAccount(addr, accountNumber, sequence, msg, baseTx)
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	base.Debug("sign transaction success")
	return txByte, nil
}

func (base *baseClient) BuildAndSendWithAccount(addr string, accountNumber, sequence uint64, msg []types.Msg, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	txByte, ctx, err := base.buildTxWithAccount(addr, accountNumber, sequence, msg, baseTx)
	if err != nil {
		return sdk.ResultTx{}, err
	}

	valid, err := base.ValidateTxSize(len(txByte), msg)
	if err != nil {
		return sdk.ResultTx{}, err
	}
	if !valid {
		base.Debug("tx is too large")
		// filter out transactions that have been sent
		// reset the maximum number of msg in each transaction
		//batch = batch / 2
		return sdk.ResultTx{}, sdk.GetError(sdk.RootCodespace, uint32(sdk.TxTooLarge))
	}
	return base.broadcastTx(txByte, ctx.Mode())
}

func (base *baseClient) BuildAndSend(msg []types.Msg, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	var res sdk.ResultTx
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
		e, ok := err.(sdk.Error)
		if ok && sdk.Code(e.Code()) == sdk.WrongSequence {
			return true
		}
		return false
	}

	onRetryFunc := func(n uint, err error) {
		_ = base.removeCache(address)
		base.Error("wrong sequence, will retry",
			"address", address, "attempts", n, "err", err.Error())
	}

	err := retry.Do(retryableFunc,
		retry.Attempts(tryThreshold),
		retry.RetryIf(retryIfFunc),
		retry.OnRetry(onRetryFunc),
		retry.LastErrorOnly(true),
	)

	if err != nil {
		return res, sdk.Wrap(err)
	}
	return res, nil
}

func (base *baseClient) prepare(baseTx sdk.BaseTx) (*sdk.Factory, error) {

	factory := sdk.NewFactory().
		WithChainID(base.cfg.ChainID).
		WithKeyManager(base.KeyManager).
		WithMode(base.cfg.Mode).
		WithSimulateAndExecute(baseTx.SimulateAndExecute).
		WithGas(base.cfg.Gas).
		WithGasAdjustment(base.cfg.GasAdjustment).
		WithSignModeHandler(base.encodingConfig.TxConfig.SignModeHandler()).
		WithTxConfig(base.encodingConfig.TxConfig).
		WithQueryFunc(base.QueryWithData).
		WithFeeGranter(base.cfg.FeeGranter).
		WithFeePayer(base.cfg.FeePayer)

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

	if !baseTx.FeeGranter.Empty() {
		factory.WithFeeGranter(baseTx.FeeGranter)
	}

	if !baseTx.FeePayer.Empty() {
		factory.WithFeePayer(baseTx.FeePayer)
	}

	if baseTx.TimeoutHeight > 0 {
		factory.WithTimeout(baseTx.TimeoutHeight)
	}

	return factory, nil
}

func (base *baseClient) prepareWithAccount(addr string, accountNumber, sequence uint64, baseTx sdk.BaseTx) (*sdk.Factory, error) {
	factory := sdk.NewFactory().
		WithChainID(base.cfg.ChainID).
		WithKeyManager(base.KeyManager).
		WithMode(base.cfg.Mode).
		WithSimulateAndExecute(baseTx.SimulateAndExecute).
		WithGas(base.cfg.Gas).
		WithSignModeHandler(base.encodingConfig.TxConfig.SignModeHandler()).
		WithTxConfig(base.encodingConfig.TxConfig).
		WithQueryFunc(base.QueryWithData).
		WithFeeGranter(base.cfg.FeeGranter).
		WithFeePayer(base.cfg.FeePayer)

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

	if !baseTx.FeeGranter.Empty() {
		factory.WithFeeGranter(baseTx.FeeGranter)
	}

	if !baseTx.FeePayer.Empty() {
		factory.WithFeePayer(baseTx.FeePayer)
	}

	if baseTx.TimeoutHeight > 0 {
		factory.WithTimeout(baseTx.TimeoutHeight)
	}

	return factory, nil
}

func (base *baseClient) ValidateTxSize(txSize int, msgs []types.Msg) (bool, sdk.Error) {
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

// NewLocker implement the function of lock, can lock resources according to conditions
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
