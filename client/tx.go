package client

import (
	"context"
	"encoding/hex"
	"strings"

	"github.com/gogo/protobuf/jsonpb"

	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/core-sdk-go/types/errors"
)

// QueryTx returns the tx info
func (base baseClient) QueryTx(hash string) (ctypes.ResultTx, error) {
	tx, err := hex.DecodeString(hash)
	if err != nil {
		return ctypes.ResultTx{}, err
	}

	res, err := base.Tx(context.Background(), tx, true)
	if err != nil {
		return ctypes.ResultTx{}, err
	}

	return *res, nil
}

func (base baseClient) QueryTxs(builder *types.EventQueryBuilder, page, size *int) (ctypes.ResultTxSearch, error) {
	query := builder.Build()
	if len(query) == 0 {
		return ctypes.ResultTxSearch{}, errors.Wrap(errors.ErrTodo, "must declare at least one tag to search")
	}

	res, err := base.TxSearch(context.Background(), query, true, page, size, "asc")
	if err != nil {
		return ctypes.ResultTxSearch{}, err
	}

	return ctypes.ResultTxSearch{
		Txs:        res.Txs,
		TotalCount: res.TotalCount,
	}, nil
}

func (base baseClient) EstimateTxGas(txBytes []byte) (uint64, error) {
	res, err := base.ABCIQuery(context.Background(), "/app/simulate", txBytes)
	if err != nil {
		return 0, err
	}

	simRes, err := parseQueryResponse(res.Response.Value)
	if err != nil {
		return 0, err
	}

	adjusted := adjustGasEstimate(simRes.GasUsed, base.cfg.GasAdjustment)
	return adjusted, nil
}

func (base *baseClient) buildTx(msgs []types.Msg, baseTx types.BaseTx) ([]byte, *types.Factory, error) {
	builder, err := base.prepare(baseTx)
	if err != nil {
		return nil, builder, errors.Wrap(errors.ErrTodo, err.Error())
	}

	txByte, err := builder.BuildAndSign(baseTx.From, msgs, false)
	if err != nil {
		return nil, builder, errors.Wrap(errors.ErrTodo, err.Error())
	}

	base.Logger().Debug("sign transaction success")
	return txByte, builder, nil
}

func (base *baseClient) buildTxWithAccount(addr string, accountNumber, sequence uint64, msgs []types.Msg, baseTx types.BaseTx) ([]byte, *types.Factory, error) {
	builder, err := base.prepareWithAccount(addr, accountNumber, sequence, baseTx)
	if err != nil {
		return nil, builder, errors.Wrap(errors.ErrTodo, err.Error())
	}

	txByte, err := builder.BuildAndSign(baseTx.From, msgs, false)
	if err != nil {
		return nil, builder, errors.Wrap(errors.ErrTodo, err.Error())
	}

	base.Logger().Debug("sign transaction success")
	return txByte, builder, nil
}

func (base baseClient) broadcastTx(txBytes []byte, mode types.BroadcastMode) (res ctypes.ResultTx, err error) {
	switch mode {
	case types.Commit:
		res, err = base.broadcastTxCommit(txBytes)
	case types.Async:
		res, err = base.broadcastTxAsync(txBytes)
	case types.Sync:
		res, err = base.broadcastTxSync(txBytes)
	default:
		err = errors.Wrapf(errors.ErrTodo, "commit mode(%s) not supported", mode)
	}
	return
}

// broadcastTxCommit broadcasts transaction bytes to a Tendermint node and waits for a commit.
func (base baseClient) broadcastTxCommit(tx []byte) (ctypes.ResultTx, error) {
	res, err := base.BroadcastTxCommit(context.Background(), tx)
	if err != nil {
		return ctypes.ResultTx{}, errors.Wrap(errors.ErrTodo, err.Error())
	}

	if !res.CheckTx.IsOK() {
		return ctypes.ResultTx{}, errors.New(res.CheckTx.Codespace, res.CheckTx.Code, res.CheckTx.Log)
	}

	if !res.DeliverTx.IsOK() {
		return ctypes.ResultTx{}, errors.New(res.DeliverTx.Codespace, res.DeliverTx.Code, res.DeliverTx.Log)
	}

	return ctypes.ResultTx{
		Hash:     res.Hash,
		Height:   res.Height,
		TxResult: res.DeliverTx,
	}, nil
}

// BroadcastTxSync broadcasts transaction bytes to a Tendermint node synchronously.
func (base baseClient) broadcastTxSync(tx []byte) (ctypes.ResultTx, error) {
	res, err := base.BroadcastTxSync(context.Background(), tx)
	if err != nil {
		return ctypes.ResultTx{}, errors.Wrap(errors.ErrTodo, err.Error())
	}

	if res.Code != 0 {
		return ctypes.ResultTx{}, errors.New(res.Codespace, res.Code, res.Log)
	}

	return ctypes.ResultTx{Hash: res.Hash}, nil
}

// BroadcastTxAsync broadcasts transaction bytes to a Tendermint node asynchronously.
func (base baseClient) broadcastTxAsync(tx []byte) (ctypes.ResultTx, error) {
	res, err := base.BroadcastTxAsync(context.Background(), tx)
	if err != nil {
		return ctypes.ResultTx{}, errors.Wrap(errors.ErrTodo, err.Error())
	}
	return ctypes.ResultTx{Hash: res.Hash}, nil
}

func adjustGasEstimate(estimate uint64, adjustment float64) uint64 {
	return uint64(adjustment * float64(estimate))
}

func parseQueryResponse(bz []byte) (types.SimulationResponse, error) {
	var simRes types.SimulationResponse
	if err := jsonpb.Unmarshal(strings.NewReader(string(bz)), &simRes); err != nil {
		return types.SimulationResponse{}, err
	}
	return simRes, nil
}
