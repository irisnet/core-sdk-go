package types

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	cryptoamino "github.com/irisnet/core-sdk-go/common/crypto/codec"
	ethsecp256k1 "github.com/irisnet/core-sdk-go/common/crypto/keys/eth_secp256k1"
	"github.com/irisnet/core-sdk-go/common/crypto/keys/secp256k1"
	"github.com/irisnet/core-sdk-go/common/crypto/keys/sm2"
	commoncryptotypes "github.com/irisnet/core-sdk-go/common/crypto/types"
	"github.com/irisnet/core-sdk-go/types/tx/signing"
	tmcrypto "github.com/tendermint/tendermint/crypto"
)

// Factory defines a client transaction factory that facilitates generating and
// signing an application-specific transaction.
type (
	// Factory implements a transaction context created in SDK modules.
	Factory struct {
		address            string
		chainID            string
		memo               string
		password           string
		timeoutHeight      uint64
		accountNumber      uint64
		sequence           uint64
		gas                uint64
		gasAdjustment      float64
		simulateAndExecute bool
		fees               Coins
		feeGranter         AccAddress
		feePayer           AccAddress
		gasPrices          DecCoins
		mode               BroadcastMode
		signMode           signing.SignMode
		signModeHandler    SignModeHandler
		keyManager         KeyManager
		txConfig           TxConfig
		queryFunc          QueryWithData
	}

	// QueryWithData implements a query method from cschain.
	QueryWithData func(string, []byte) ([]byte, int64, error)
)

// NewFactory return a point of the instance of Factory.
func NewFactory() *Factory {
	return &Factory{}
}

// ChainID returns the chainID of the current chain.
func (f *Factory) ChainID() string { return f.chainID }

// Gas returns the gas of the transaction.
func (f *Factory) Gas() uint64 { return f.gas }

// GasAdjustment returns the gasAdjustment.
func (f Factory) GasAdjustment() float64 { return f.gasAdjustment }

// Fees returns the fee of the transaction.
func (f *Factory) Fees() Coins { return f.fees }

// Sequence returns the sequence of the account.
func (f *Factory) Sequence() uint64 { return f.sequence }

// Memo returns memo.
func (f *Factory) Memo() string { return f.memo }

// AccountNumber returns accountNumber.
func (f *Factory) AccountNumber() uint64 { return f.accountNumber }

// KeyManager returns keyManager.
func (f *Factory) KeyManager() KeyManager { return f.keyManager }

// Mode returns mode.
func (f *Factory) Mode() BroadcastMode { return f.mode }

// SimulateAndExecute returns the option to simulateAndExecute and then execute the transaction
// using the gas from the simulation results
func (f *Factory) SimulateAndExecute() bool { return f.simulateAndExecute }

// Password returns password.
func (f *Factory) Password() string { return f.password }

// Address returns the address.
func (f *Factory) Address() string { return f.address }

// WithChainID returns a pointer of the context with an updated ChainID.
func (f *Factory) WithChainID(chainID string) *Factory {
	f.chainID = chainID
	return f
}

// WithGas returns a pointer of the context with an updated Gas.
func (f *Factory) WithGas(gas uint64) *Factory {
	f.gas = gas
	return f
}

// WithGasAdjustment returns a pointer of the context with an updated gasAdjustment.
func (f *Factory) WithGasAdjustment(gasAdjustment float64) *Factory {
	f.gasAdjustment = gasAdjustment
	return f
}

// WithFee returns a pointer of the context with an updated Fee.
func (f *Factory) WithFee(fee Coins) *Factory {
	f.fees = fee
	return f
}

// WithFeeGranter returns a pointer of the context with an updated FeeGranter.
func (f *Factory) WithFeeGranter(feeGranter AccAddress) *Factory {
	f.feeGranter = feeGranter
	return f
}

// WithFeePayer returns a pointer of the context with an updated FeePayer.
func (f *Factory) WithFeePayer(feePayer AccAddress) *Factory {
	f.feePayer = feePayer
	return f
}

// WithSequence returns a pointer of the context with an updated sequence number.
func (f *Factory) WithSequence(sequence uint64) *Factory {
	f.sequence = sequence
	return f
}

// WithMemo returns a pointer of the context with an updated memo.
func (f *Factory) WithMemo(memo string) *Factory {
	f.memo = memo
	return f
}

// WithAccountNumber returns a pointer of the context with an account number.
func (f *Factory) WithAccountNumber(accnum uint64) *Factory {
	f.accountNumber = accnum
	return f
}

// WithKeyManager returns a pointer of the context with a KeyManager.
func (f *Factory) WithKeyManager(keyManager KeyManager) *Factory {
	f.keyManager = keyManager
	return f
}

// WithMode returns a pointer of the context with a Mode.
func (f *Factory) WithMode(mode BroadcastMode) *Factory {
	f.mode = mode
	return f
}

// WithSimulateAndExecute returns a pointer of the context with a simulateAndExecute.
func (f *Factory) WithSimulateAndExecute(simulate bool) *Factory {
	f.simulateAndExecute = simulate
	return f
}

// WithPassword returns a pointer of the context with a password.
func (f *Factory) WithPassword(password string) *Factory {
	f.password = password
	return f
}

// WithAddress returns a pointer of the context with a password.
func (f *Factory) WithAddress(address string) *Factory {
	f.address = address
	return f
}

// WithTxConfig returns a pointer of the context with an TxConfig
func (f *Factory) WithTxConfig(txConfig TxConfig) *Factory {
	f.txConfig = txConfig
	return f
}

// WithSignModeHandler returns a pointer of the context with an signModeHandler.
func (f *Factory) WithSignModeHandler(signModeHandler SignModeHandler) *Factory {
	f.signModeHandler = signModeHandler
	return f
}

// WithQueryFunc returns a pointer of the context with an queryFunc.
func (f *Factory) WithQueryFunc(queryFunc QueryWithData) *Factory {
	f.queryFunc = queryFunc
	return f
}

// WithTimeout Timeout for accessing the blockchain (such as query transactions, broadcast transactions, etc.)
func (f *Factory) WithTimeout(height uint64) *Factory {
	f.timeoutHeight = height
	return f
}

func (f *Factory) BuildAndSign(name string, msgs []Msg, json bool) ([]byte, error) {
	tx, err := f.BuildUnsignedTx(msgs)
	if err != nil {
		return nil, err
	}

	if err = f.Sign(name, tx); err != nil {
		return nil, err
	}

	if json {
		txBytes, err := f.txConfig.TxJSONEncoder()(tx.GetTx())
		if err != nil {
			return nil, err
		}
		return txBytes, nil
	}

	txBytes, err := f.txConfig.TxEncoder()(tx.GetTx())
	if err != nil {
		return nil, err
	}

	return txBytes, nil
}

func (f *Factory) BuildTxWithoutKeyDao(pubkey []byte, algo string, msgs []Msg) ([]byte, error) {
	tx, err := f.BuildUnsignedTx(msgs)
	if err != nil {
		return []byte{}, err
	}

	signMode := f.signMode
	if signMode == signing.SignMode_SIGN_MODE_UNSPECIFIED {
		// use the SignModeHandler's default mode if unspecified
		signMode = f.txConfig.SignModeHandler().DefaultMode()
	}
	signerData := SignerData{
		ChainID:       f.chainID,
		AccountNumber: f.accountNumber,
		Sequence:      f.sequence,
	}

	pubKey, err := cryptoamino.PubKeyFromBytes(pubkey)
	if err != nil {
		return nil, nil
	}

	publicKey := FromTmPubKey(algo, pubKey)

	sigData := signing.SingleSignatureData{
		SignMode:  signMode,
		Signature: nil,
	}
	sig := signing.SignatureV2{
		PubKey:   publicKey,
		Data:     &sigData,
		Sequence: f.Sequence(),
	}
	if err := tx.SetSignatures(sig); err != nil {
		return []byte{}, err
	}

	// Generate the bytes to be signed.
	signBytes, err := f.signModeHandler.GetSignBytes(signMode, signerData, tx.GetTx())
	if err != nil {
		return []byte{}, err
	}

	return crypto.Keccak256Hash(signBytes).Bytes(), nil
}

func (f *Factory) SetUnsignedTxSignature(pubkey []byte, algo string, msgs []Msg, signedData []byte) ([]byte, error) {
	tx, err := f.BuildUnsignedTx(msgs)
	if err != nil {
		return []byte{}, err
	}

	signMode := f.signMode
	if signMode == signing.SignMode_SIGN_MODE_UNSPECIFIED {
		// use the SignModeHandler's default mode if unspecified
		signMode = f.txConfig.SignModeHandler().DefaultMode()
	}

	pubKey, err := cryptoamino.PubKeyFromBytes(pubkey)
	if err != nil {
		return nil, nil
	}

	publicKey := FromTmPubKey(algo, pubKey)

	// Construct the SignatureV2 struct
	sigData := signing.SingleSignatureData{
		SignMode:  signMode,
		Signature: signedData,
	}
	sig := signing.SignatureV2{
		PubKey:   publicKey,
		Data:     &sigData,
		Sequence: f.Sequence(),
	}

	// And here the tx is populated with the signature
	tx.SetSignatures(sig)

	txBytes, err := f.txConfig.TxEncoder()(tx.GetTx())
	if err != nil {
		return nil, err
	}

	return txBytes, nil
}

func (f *Factory) BuildUnsignedTx(msgs []Msg) (TxBuilder, error) {
	if f.chainID == "" {
		return nil, fmt.Errorf("chain ID required but not specified")
	}

	fees := f.fees

	if !f.gasPrices.IsZero() {
		if !fees.IsZero() {
			return nil, errors.New("cannot provide both fees and gas prices")
		}

		glDec := NewDec(int64(f.gas))

		// Derive the fees based on the provided gas prices, where
		// fee = ceil(gasPrice * gasLimit).
		fees = make(Coins, len(f.gasPrices))

		for i, gp := range f.gasPrices {
			fee := gp.Amount.Mul(glDec)
			fees[i] = NewCoin(gp.Denom, fee.Ceil().RoundInt())
		}
	}

	tx := f.txConfig.NewTxBuilder()

	if err := tx.SetMsgs(msgs...); err != nil {
		return nil, err
	}

	tx.SetMemo(f.memo)
	tx.SetFeeAmount(fees)
	tx.SetGasLimit(f.gas)
	tx.SetFeeGranter(f.feeGranter)
	tx.SetFeePayer(f.feePayer)
	tx.SetTimeoutHeight(f.timeoutHeight)
	//f.txBuilder.SetTimeoutHeight(f.TimeoutHeight())

	return tx, nil
}

// Sign signs a transaction given a name, passphrase, and a single message to
// signed. An error is returned if signing fails.
func (f *Factory) Sign(name string, txBuilder TxBuilder) error {
	signMode := f.signMode
	if signMode == signing.SignMode_SIGN_MODE_UNSPECIFIED {
		// use the SignModeHandler's default mode if unspecified
		signMode = f.txConfig.SignModeHandler().DefaultMode()
	}
	signerData := SignerData{
		ChainID:       f.chainID,
		AccountNumber: f.accountNumber,
		Sequence:      f.sequence,
	}

	pubkey, _, err := f.keyManager.Find(name, f.password)
	if err != nil {
		return err
	}

	// For SIGN_MODE_DIRECT, calling SetSignatures calls setSignerInfos on
	// Factory under the hood, and SignerInfos is needed to generated the
	// sign bytes. This is the reason for setting SetSignatures here, with a
	// nil signature.
	//
	// Note: this line is not needed for SIGN_MODE_LEGACY_AMINO, but putting it
	// also doesn't affect its generated sign bytes, so for code's simplicity
	// sake, we put it here.
	sigData := signing.SingleSignatureData{
		SignMode:  signMode,
		Signature: nil,
	}
	sig := signing.SignatureV2{
		PubKey:   pubkey,
		Data:     &sigData,
		Sequence: f.Sequence(),
	}
	if err := txBuilder.SetSignatures(sig); err != nil {
		return err
	}

	// Generate the bytes to be signed.
	signBytes, err := f.signModeHandler.GetSignBytes(signMode, signerData, txBuilder.GetTx())
	if err != nil {
		return err
	}

	// Sign those bytes
	sigBytes, _, err := f.keyManager.Sign(name, f.password, signBytes)
	if err != nil {
		return err
	}

	// Construct the SignatureV2 struct
	sigData = signing.SingleSignatureData{
		SignMode:  signMode,
		Signature: sigBytes,
	}
	sig = signing.SignatureV2{
		PubKey:   pubkey,
		Data:     &sigData,
		Sequence: f.Sequence(),
	}

	// And here the tx is populated with the signature
	return txBuilder.SetSignatures(sig)
}

func FromTmPubKey(Algo string, pubKey tmcrypto.PubKey) commoncryptotypes.PubKey {
	var pubkey commoncryptotypes.PubKey
	pubkeyBytes := pubKey.Bytes()
	switch Algo {
	case "sm2":
		pubkey = &sm2.PubKey{Key: pubkeyBytes}
	case "secp256k1":
		pubkey = &secp256k1.PubKey{Key: pubkeyBytes}
	case "eth_secp256k1":
		pubkey = &ethsecp256k1.PubKey{Key: pubkeyBytes}
	}
	return pubkey
}
