package types

import (
	"encoding/base64"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/encoding"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

type Block struct {
	tmtypes.Header `json:"header"`
	Data           `json:"data"`
	Evidence       tmtypes.EvidenceData `json:"evidence"`
	LastCommit     *tmtypes.Commit      `json:"last_commit"`
}

type Data struct {
	Txs []StdTx `json:"txs"`
}

type BlockResult struct {
	Height  int64         `json:"height"`
	Results ABCIResponses `json:"results"`
}

type BlockDetail struct {
	BlockID     tmtypes.BlockID `json:"block_id"`
	Block       Block           `json:"block"`
	BlockResult BlockResult     `json:"block_result"`
}

type ABCIResponses struct {
	DeliverTx  []TxResult
	EndBlock   ResultEndBlock
	BeginBlock ResultBeginBlock
}

type ResultBeginBlock struct {
	Events StringEvents `json:"events"`
}

type ResultEndBlock struct {
	Events           StringEvents      `json:"events"`
	ValidatorUpdates []ValidatorUpdate `json:"validator_updates"`
}

func ParseValidatorUpdate(updates []abci.ValidatorUpdate) []ValidatorUpdate {
	var vUpdates []ValidatorUpdate
	for _, v := range updates {
		pubkey, _ := encoding.PubKeyFromProto(v.PubKey)
		vUpdates = append(
			vUpdates,
			ValidatorUpdate{
				PubKey: PubKey{
					Type:  pubkey.Type(),
					Value: base64.StdEncoding.EncodeToString(pubkey.Bytes()),
				},
				Power: v.Power,
			},
		)
	}
	return vUpdates
}

func ParseBlockResult(res *ctypes.ResultBlockResults) BlockResult {
	var txResults = make([]TxResult, len(res.TxsResults))
	for i, r := range res.TxsResults {
		txResults[i] = TxResult{
			Code:      r.Code,
			Log:       r.Log,
			GasWanted: r.GasWanted,
			GasUsed:   r.GasUsed,
			Events:    StringifyEvents(r.Events),
		}
	}
	return BlockResult{
		Height: res.Height,
		Results: ABCIResponses{
			DeliverTx: txResults,
			EndBlock: ResultEndBlock{
				Events:           StringifyEvents(res.EndBlockEvents),
				ValidatorUpdates: ParseValidatorUpdate(res.ValidatorUpdates),
			},
			BeginBlock: ResultBeginBlock{
				Events: StringifyEvents(res.BeginBlockEvents),
			},
		},
	}
}
