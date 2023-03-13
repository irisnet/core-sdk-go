package ethermint

type IClient interface {
	BuildEvmTx(hexData string, feePayerAddr string, evmDemon string) error
}
