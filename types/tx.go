package types

const (
	// maxMemoCharacters = 100
	// txSigLimit        = 7
	// maxGasWanted = uint64((1 << 63) - 1)

	Sync   BroadcastMode = "sync"
	Async  BroadcastMode = "async"
	Commit BroadcastMode = "commit"
)

type BroadcastMode string

type Msgs []Msg

func (m Msgs) Len() int {
	return len(m)
}

func (m Msgs) Sub(begin, end int) SplitAble {
	return m[begin:end]
}

type BaseTx struct {
	From               string        `json:"from"`
	Password           string        `json:"password"`
	Gas                uint64        `json:"gas"`
	Fee                DecCoins      `json:"fee"`
	Memo               string        `json:"memo"`
	Mode               BroadcastMode `json:"broadcast_mode"`
	SimulateAndExecute bool          `json:"simulate_and_execute"`
	GasAdjustment      float64       `json:"gas_adjustment"`
}
