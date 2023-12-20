package types

// BaseAccount defines the basic structure of the account
type BaseAccount struct {
	Address       string `json:"address"`
	PubKey        string `json:"public_key"`
	PubKeyType    string `json:"pubkey_type"`
	AccountNumber uint64 `json:"account_number"`
	Sequence      uint64 `json:"sequence"`
}
