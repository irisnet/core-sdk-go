package types

import (
	"errors"
	"math/big"

	sdkcodectypes "github.com/irisnet/core-sdk-go/common/codec/types"
	sdktypes "github.com/irisnet/core-sdk-go/types"
	sdktxtypes "github.com/irisnet/core-sdk-go/types/tx"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/irisnet/core-sdk-go/ethermint/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

var (
	_ sdktypes.Msg = &MsgEthereumTx{}
	_ sdktypes.Tx  = &MsgEthereumTx{}
	//_ ante.GasTx   = &MsgEthereumTx{}

	//_ codectypes.UnpackInterfacesMessage = MsgEthereumTx{}
)

// message type and route constants
const (
	// TypeMsgEthereumTx defines the type string of an Ethereum transaction
	TypeMsgEthereumTx = "ethereum_tx"
)

// NewTx returns a reference to a new Ethereum transaction message.
func NewTx(
	chainID *big.Int, nonce uint64, to *common.Address, amount *big.Int,
	gasLimit uint64, gasPrice, gasFeeCap, gasTipCap *big.Int, input []byte, accesses *ethtypes.AccessList,
) *MsgEthereumTx {
	return newMsgEthereumTx(chainID, nonce, to, amount, gasLimit, gasPrice, gasFeeCap, gasTipCap, input, accesses)
}

func newMsgEthereumTx(
	chainID *big.Int, nonce uint64, to *common.Address, amount *big.Int,
	gasLimit uint64, gasPrice, gasFeeCap, gasTipCap *big.Int, input []byte, accesses *ethtypes.AccessList,
) *MsgEthereumTx {
	var (
		cid, amt, gp *sdk.Int
		toAddr       string
		txData       TxData
	)

	if to != nil {
		toAddr = to.Hex()
	}

	if amount != nil {
		amountInt := sdk.NewIntFromBigInt(amount)
		amt = &amountInt
	}

	if chainID != nil {
		chainIDInt := sdk.NewIntFromBigInt(chainID)
		cid = &chainIDInt
	}

	if gasPrice != nil {
		gasPriceInt := sdk.NewIntFromBigInt(gasPrice)
		gp = &gasPriceInt
	}

	switch {
	case accesses == nil:
		txData = &LegacyTx{
			Nonce:    nonce,
			To:       toAddr,
			Amount:   amt,
			GasLimit: gasLimit,
			GasPrice: gp,
			Data:     input,
		}
	case accesses != nil && gasFeeCap != nil && gasTipCap != nil:
		gtc := sdk.NewIntFromBigInt(gasTipCap)
		gfc := sdk.NewIntFromBigInt(gasFeeCap)

		txData = &DynamicFeeTx{
			ChainID:   cid,
			Nonce:     nonce,
			To:        toAddr,
			Amount:    amt,
			GasLimit:  gasLimit,
			GasTipCap: &gtc,
			GasFeeCap: &gfc,
			Data:      input,
			Accesses:  NewAccessList(accesses),
		}
	case accesses != nil:
		txData = &AccessListTx{
			ChainID:  cid,
			Nonce:    nonce,
			To:       toAddr,
			Amount:   amt,
			GasLimit: gasLimit,
			GasPrice: gp,
			Data:     input,
			Accesses: NewAccessList(accesses),
		}
	default:
	}

	dataAny, err := PackTxData(txData)
	if err != nil {
		panic(err)
	}

	return &MsgEthereumTx{Data: dataAny}
}

// fromEthereumTx populates the message fields from the given ethereum transaction
func (msg *MsgEthereumTx) FromEthereumTx(tx *ethtypes.Transaction) error {
	txData, err := NewTxDataFromTx(tx)
	if err != nil {
		return err
	}

	anyTxData, err := PackTxData(txData)
	if err != nil {
		return err
	}

	msg.Data = anyTxData
	msg.Size_ = float64(tx.Size())
	msg.Hash = tx.Hash().Hex()
	return nil
}

// Route returns the route value of an MsgEthereumTx.
func (msg MsgEthereumTx) Route() string { return RouterKey }

// Type returns the type value of an MsgEthereumTx.
func (msg MsgEthereumTx) Type() string { return TypeMsgEthereumTx }

// ValidateBasic implements the sdktypes.Msg interface. It performs basic validation
// checks of a Transaction. If returns an error if validation fails.
func (msg MsgEthereumTx) ValidateBasic() error {
	if msg.From != "" {
		if err := types.ValidateAddress(msg.From); err != nil {
			return sdkerrors.Wrap(err, "invalid from address")
		}
	}

	txData, err := UnpackTxData(msg.Data)
	if err != nil {
		return sdkerrors.Wrap(err, "failed to unpack tx data")
	}

	return txData.Validate()
}

// GetMsgs returns a single MsgEthereumTx as an sdktypes.Msg.
func (msg *MsgEthereumTx) GetMsgs() []sdktypes.Msg {
	return []sdktypes.Msg{msg}
}

// GetSigners returns the expected signers for an Ethereum transaction message.
// For such a message, there should exist only a single 'signer'.
//
// NOTE: This method panics if 'Sign' hasn't been called first.
func (msg *MsgEthereumTx) GetSigners() []sdktypes.AccAddress {
	data, err := UnpackTxData(msg.Data)
	if err != nil {
		panic(err)
	}

	sender, err := msg.GetSender(data.GetChainID())
	if err != nil {
		panic(err)
	}

	signer := sdktypes.AccAddress(sender.Bytes())
	return []sdktypes.AccAddress{signer}
}

// GetSignBytes returns the Amino bytes of an Ethereum transaction message used
// for signing.
//
// NOTE: This method cannot be used as a chain ID is needed to create valid bytes
// to sign over. Use 'RLPSignBytes' instead.
func (msg MsgEthereumTx) GetSignBytes() []byte {
	panic("must use 'RLPSignBytes' with a chain ID to get the valid bytes to sign")
}

// GetGas implements the GasTx interface. It returns the GasLimit of the transaction.
func (msg MsgEthereumTx) GetGas() uint64 {
	txData, err := UnpackTxData(msg.Data)
	if err != nil {
		return 0
	}
	return txData.GetGas()
}

// GetFee returns the fee for non dynamic fee tx
func (msg MsgEthereumTx) GetFee() *big.Int {
	txData, err := UnpackTxData(msg.Data)
	if err != nil {
		return nil
	}
	return txData.Fee()
}

// GetEffectiveFee returns the fee for dynamic fee tx
func (msg MsgEthereumTx) GetEffectiveFee(baseFee *big.Int) *big.Int {
	txData, err := UnpackTxData(msg.Data)
	if err != nil {
		return nil
	}
	return txData.EffectiveFee(baseFee)
}

// GetFrom loads the ethereum sender address from the sigcache and returns an
// sdktypes.AccAddress from its bytes
func (msg *MsgEthereumTx) GetFrom() sdktypes.AccAddress {
	if msg.From == "" {
		return nil
	}

	return common.HexToAddress(msg.From).Bytes()
}

// GetFeePayer loads the ethereum sender address from the sigcache and returns an
// sdktypes.AccAddress from its bytes
func (msg *MsgEthereumTx) GetFeePayer() sdktypes.AccAddress {
	if msg.FeePayer == "" {
		return nil
	}

	return common.HexToAddress(msg.FeePayer).Bytes()
}

func (msg *MsgEthereumTx) SetFeePayer(feePayer string) {
	msg.FeePayer = feePayer
}

// AsTransaction creates an Ethereum Transaction type from the msg fields
func (msg MsgEthereumTx) AsTransaction() *ethtypes.Transaction {
	txData, err := UnpackTxData(msg.Data)
	if err != nil {
		return nil
	}

	return ethtypes.NewTx(txData.AsEthereumData())
}

// AsMessage creates an Ethereum core.Message from the msg fields
func (msg MsgEthereumTx) AsMessage(signer ethtypes.Signer, baseFee *big.Int) (core.Message, error) {
	return msg.AsTransaction().AsMessage(signer, baseFee)
}

// GetSender extracts the sender address from the signature values using the latest signer for the given chainID.
func (msg *MsgEthereumTx) GetSender(chainID *big.Int) (common.Address, error) {

	signer := ethtypes.LatestSignerForChainID(chainID)
	from, err := signer.Sender(msg.AsTransaction())
	if err != nil {
		return common.Address{}, err
	}

	msg.From = from.Hex()
	return from, nil
}

//// UnpackInterfaces implements UnpackInterfacesMesssage.UnpackInterfaces
//func (msg MsgEthereumTx) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
//	return unpacker.UnpackAny(msg.Data, new(TxData))
//}

// UnmarshalBinary decodes the canonical encoding of transactions.
func (msg *MsgEthereumTx) UnmarshalBinary(b []byte) error {
	tx := &ethtypes.Transaction{}
	if err := tx.UnmarshalBinary(b); err != nil {
		return err
	}
	return msg.FromEthereumTx(tx)
}

// BuildTx builds the canonical cosmos tx from ethereum msg
func (msg *MsgEthereumTx) BuildTx(b sdktypes.TxBuilder, evmDenom string) (sdktypes.Tx, error) {
	builder, ok := b.(sdktxtypes.ExtensionOptionsTxBuilder)
	if !ok {
		return nil, errors.New("unsupported builder")
	}
	option, err := sdkcodectypes.NewAnyWithValue(&ExtensionOptionsEthereumTx{})
	if err != nil {
		return nil, err
	}

	txData, err := UnpackTxData(msg.Data)
	if err != nil {
		return nil, err
	}
	fees := make(sdktypes.Coins, 0)
	feeAmt := sdktypes.NewIntFromBigInt(txData.Fee())
	if feeAmt.Sign() > 0 {
		fees = append(fees, sdktypes.NewCoin(evmDenom, feeAmt))
	}

	builder.SetExtensionOptions(option)
	err = builder.SetMsgs(msg)
	if err != nil {
		return nil, err
	}
	builder.SetFeeAmount(fees)
	builder.SetGasLimit(msg.GetGas())
	tx := builder.GetTx()
	return tx, nil
}
