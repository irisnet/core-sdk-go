package types

import "github.com/cosmos/cosmos-sdk/types"

const (
	// maxMemoCharacters = 100
	// txSigLimit        = 7
	maxGasWanted = uint64((1 << 63) - 1)

	Sync   BroadcastMode = "sync"
	Async  BroadcastMode = "async"
	Commit BroadcastMode = "commit"
)

type (
	BroadcastMode string
)

type BaseTx struct {
	From               string           `json:"from"`
	Password           string           `json:"password"`
	Gas                uint64           `json:"gas"`
	Fee                types.DecCoins   `json:"fee"`
	FeePayer           types.AccAddress `json:"fee_payer"`
	FeeGranter         types.AccAddress `json:"fee_granter"`
	Memo               string           `json:"memo"`
	Mode               BroadcastMode    `json:"broadcast_mode"`
	SimulateAndExecute bool             `json:"simulate_and_execute"`
	GasAdjustment      float64          `json:"gas_adjustment"`
	TimeoutHeight      uint64           `json:"timeout_height"`
}

// ResultTx encapsulates the return result of the transaction. When the transaction fails,
// it is an empty object. The specific error information can be obtained through the Error interface.
type ResultTx struct {
	GasWanted int64              `json:"gas_wanted"`
	GasUsed   int64              `json:"gas_used"`
	Data      []byte             `json:"data"`
	Events    types.StringEvents `json:"events"`
	Hash      string             `json:"hash"`
	Height    int64              `json:"height"`
}
