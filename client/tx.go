package client

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/jsonpb"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	sdk "github.com/irisnet/core-sdk-go/types"
)

// QueryTx returns the tx info
func (base *baseClient) QueryTx(hash string) (*types.TxResponse, error) {
	tx, err := hex.DecodeString(hash)
	if err != nil {
		return nil, err
	}

	res, err := base.Tx(context.Background(), tx, true)
	if err != nil {
		return nil, err
	}

	resBlocks, err := base.getResultBlocks([]*ctypes.ResultTx{res})
	if err != nil {
		return nil, err
	}
	return base.mkTxResult(res, resBlocks[res.Height])
}

func (base *baseClient) QueryTxs(events []string, page, limit int, orderBy string) (*types.SearchTxsResult, error) {
	if len(events) == 0 {
		return nil, errors.New("must declare at least one tag to search")
	}
	if page <= 0 {
		return nil, errors.New("page must be greater than 0")
	}

	if limit <= 0 {
		return nil, errors.New("limit must be greater than 0")
	}

	// XXX: implement ANY
	query := strings.Join(events, " AND ")

	// TODO: this may not always need to be proven
	// https://github.com/cosmos/cosmos-sdk/issues/6807
	resTxs, err := base.TxSearch(context.Background(), query, true, &page, &limit, orderBy)
	if err != nil {
		return nil, err
	}
	resBlocks, err := base.getResultBlocks(resTxs.Txs)
	if err != nil {
		return nil, err
	}

	txs, err := base.formatTxResults(resTxs.Txs, resBlocks)
	if err != nil {
		return nil, err
	}

	result := types.NewSearchTxsResult(uint64(resTxs.TotalCount), uint64(len(txs)), uint64(page), uint64(limit), txs)

	return result, nil

}

func (base *baseClient) BlockMetadata(height int64) (sdk.BlockDetailMetadata, error) {
	block, err := base.Block(context.Background(), &height)
	if err != nil {
		return sdk.BlockDetailMetadata{}, err
	}
	blockMetadata, err := json.Marshal(block)
	if err != nil {
		return sdk.BlockDetailMetadata{}, err
	}

	blockResult, err := base.BlockResults(context.Background(), &height)
	if err != nil {
		return sdk.BlockDetailMetadata{}, err
	}
	blockResultMetadata, err := json.Marshal(blockResult)
	if err != nil {
		return sdk.BlockDetailMetadata{}, err
	}

	return sdk.BlockDetailMetadata{
		Block:       blockMetadata,
		BlockResult: blockResultMetadata,
	}, nil
}

func (base *baseClient) EstimateTxGas(txBytes []byte) (uint64, error) {
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

func (base *baseClient) buildTx(msgs []types.Msg, baseTx sdk.BaseTx) ([]byte, *sdk.Factory, sdk.Error) {
	builder, err := base.prepare(baseTx)
	if err != nil {
		return nil, builder, sdk.Wrap(err)
	}
	txByte, err := builder.BuildAndSign(baseTx.From, msgs, false)
	if err != nil {
		return nil, builder, sdk.Wrap(err)
	}
	base.Debug("sign transaction success")
	return txByte, builder, nil
}

func (base *baseClient) buildTxWithAccount(addr string, accountNumber, sequence uint64, msgs []types.Msg, baseTx sdk.BaseTx) ([]byte, *sdk.Factory, sdk.Error) {
	builder, err := base.prepareWithAccount(addr, accountNumber, sequence, baseTx)
	if err != nil {
		return nil, builder, sdk.Wrap(err)
	}

	txByte, err := builder.BuildAndSign(baseTx.From, msgs, false)
	if err != nil {
		return nil, builder, sdk.Wrap(err)
	}

	base.Debug("sign transaction success")
	return txByte, builder, nil
}

func (base *baseClient) broadcastTx(txBytes []byte, mode sdk.BroadcastMode) (res sdk.ResultTx, err sdk.Error) {
	switch mode {
	case sdk.Commit:
		res, err = base.broadcastTxCommit(txBytes)
	case sdk.Async:
		res, err = base.broadcastTxAsync(txBytes)
	case sdk.Sync:
		res, err = base.broadcastTxSync(txBytes)
	default:
		err = sdk.Wrapf("commit mode(%s) not supported", mode)
	}
	return
}

// broadcastTxCommit broadcasts transaction bytes to a Tendermint node
// and waits for a commit.
func (base *baseClient) broadcastTxCommit(tx []byte) (sdk.ResultTx, sdk.Error) {
	res, err := base.BroadcastTxCommit(context.Background(), tx)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	if !res.CheckTx.IsOK() {
		return sdk.ResultTx{}, sdk.GetError(res.CheckTx.Codespace, res.CheckTx.Code, res.CheckTx.Log)
	}

	if !res.DeliverTx.IsOK() {
		return sdk.ResultTx{}, sdk.GetError(res.DeliverTx.Codespace, res.DeliverTx.Code, res.DeliverTx.Log)
	}

	return sdk.ResultTx{
		GasWanted: res.DeliverTx.GasWanted,
		GasUsed:   res.DeliverTx.GasUsed,
		Events:    types.StringifyEvents(res.DeliverTx.Events),
		Hash:      res.Hash.String(),
		Height:    res.Height,
	}, nil
}

// BroadcastTxSync broadcasts transaction bytes to a Tendermint node
// synchronously.
func (base *baseClient) broadcastTxSync(tx []byte) (sdk.ResultTx, sdk.Error) {
	res, err := base.BroadcastTxSync(context.Background(), tx)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	if res.Code != 0 {
		return sdk.ResultTx{}, sdk.GetError(res.Codespace, res.Code, res.Log)
	}

	return sdk.ResultTx{Hash: res.Hash.String()}, nil
}

// BroadcastTxAsync broadcasts transaction bytes to a Tendermint node
// asynchronously.
func (base *baseClient) broadcastTxAsync(tx []byte) (sdk.ResultTx, sdk.Error) {
	res, err := base.BroadcastTxAsync(context.Background(), tx)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	return sdk.ResultTx{Hash: res.Hash.String()}, nil
}

func (base *baseClient) getResultBlocks(resTxs []*ctypes.ResultTx) (map[int64]*ctypes.ResultBlock, error) {
	resBlocks := make(map[int64]*ctypes.ResultBlock)
	for _, resTx := range resTxs {
		if _, ok := resBlocks[resTx.Height]; !ok {
			resBlock, err := base.Block(context.Background(), &resTx.Height)
			if err != nil {
				return nil, err
			}

			resBlocks[resTx.Height] = resBlock
		}
	}

	return resBlocks, nil
}

func (base *baseClient) mkTxResult(resTx *ctypes.ResultTx, resBlock *ctypes.ResultBlock) (*types.TxResponse, error) {
	txb, err := base.encodingConfig.TxConfig.TxDecoder()(resTx.Tx)
	if err != nil {
		return nil, err
	}

	p, ok := txb.(intoAny)
	if !ok {
		return nil, fmt.Errorf("expecting a type implementing intoAny, got: %T", txb)
	}
	any := p.AsAny()

	return types.NewResponseResultTx(resTx, any, resBlock.Block.Time.Format(time.RFC3339)), nil
}

// formatTxResults parses the indexed txs into a slice of TxResponse objects.
func (base *baseClient) formatTxResults(resTxs []*ctypes.ResultTx, resBlocks map[int64]*ctypes.ResultBlock) ([]*types.TxResponse, error) {
	var err error
	out := make([]*types.TxResponse, len(resTxs))
	for i := range resTxs {
		out[i], err = base.mkTxResult(resTxs[i], resBlocks[resTxs[i].Height])
		if err != nil {
			return nil, err
		}
	}

	return out, nil
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

// Deprecated: this interface is used only internally for scenario we are
// deprecating (StdTxConfig support)
type intoAny interface {
	AsAny() *codectypes.Any
}
