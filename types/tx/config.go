package tx

import (
	"fmt"

	"github.com/irisnet/core-sdk-go/codec"
	"github.com/irisnet/core-sdk-go/types"
	signingtypes "github.com/irisnet/core-sdk-go/types/tx/signing"
)

type config struct {
	handler     types.SignModeHandler
	decoder     types.TxDecoder
	encoder     types.TxEncoder
	jsonDecoder types.TxDecoder
	jsonEncoder types.TxEncoder
	protoCodec  *codec.ProtoCodec
}

// NewTxConfig returns a new protobuf TxConfig using the provided ProtoCodec and sign modes. The
// first enabled sign mode will become the default sign mode.
func NewTxConfig(protoCodec *codec.ProtoCodec, enabledSignModes []signingtypes.SignMode) types.TxConfig {
	return &config{
		handler:     MakeSignModeHandler(enabledSignModes),
		decoder:     DefaultTxDecoder(protoCodec),
		encoder:     DefaultTxEncoder(),
		jsonDecoder: DefaultJSONTxDecoder(protoCodec),
		jsonEncoder: DefaultJSONTxEncoder(),
		protoCodec:  protoCodec,
	}
}

func (g config) NewTxBuilder() types.TxBuilder {
	return newBuilder()
}

// WrapTxBuilder returns a builder from provided transaction
func (g config) WrapTxBuilder(newTx types.Tx) (types.TxBuilder, error) {
	newBuilder, ok := newTx.(*wrapper)
	if !ok {
		return nil, fmt.Errorf("expected %T, got %T", &wrapper{}, newTx)
	}

	return newBuilder, nil
}

func (g config) SignModeHandler() types.SignModeHandler {
	return g.handler
}

func (g config) TxEncoder() types.TxEncoder {
	return g.encoder
}

func (g config) TxDecoder() types.TxDecoder {
	return g.decoder
}

func (g config) TxJSONEncoder() types.TxEncoder {
	return g.jsonEncoder
}

func (g config) TxJSONDecoder() types.TxDecoder {
	return g.jsonDecoder
}
