package types

type TxResult struct {
	Code      uint32       `json:"code"`
	Log       string       `json:"log"`
	GasWanted int64        `json:"gas_wanted"`
	GasUsed   int64        `json:"gas_used"`
	Events    StringEvents `json:"events"`
}
