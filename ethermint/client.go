package ethermint

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	codectypes "github.com/irisnet/core-sdk-go/common/codec/types"
	evmtypes "github.com/irisnet/core-sdk-go/ethermint/x/evm/types"
	sdktypes "github.com/irisnet/core-sdk-go/types"
	"github.com/pkg/errors"
)

type Client struct {
	sdktypes.BaseClient
	txConfig sdktypes.TxConfig
}

// NewClient grant NewClient
func NewClient(bc sdktypes.BaseClient, txConfig sdktypes.TxConfig) Client {
	return Client{
		BaseClient: bc,
		txConfig:   txConfig,
	}
}

func (cli *Client) RegisterInterfaceTypes(registry codectypes.InterfaceRegistry) {
	evmtypes.RegisterInterfaces(registry)
}

func (cli *Client) BuildEvmTx(hexData string, feePayerAddr string, evmDemon string) ([]byte, error) {
	data, err := hexutil.Decode(hexData)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode ethereum tx hex bytes")
	}
	msg := &evmtypes.MsgEthereumTx{}
	if err := msg.UnmarshalBinary(data); err != nil {
		return nil, err
	}
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	if len(feePayerAddr) > 0 {
		msg.SetFeePayer(feePayerAddr)
	}

	builder := cli.txConfig.NewTxBuilder()
	tx, err := msg.BuildTx(builder, evmDemon)
	if err != nil {
		return nil, err
	}

	return cli.txConfig.TxEncoder()(tx)
}
