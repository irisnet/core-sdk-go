package ethermint

import (
	"fmt"
	"testing"

	sdk "github.com/irisnet/core-sdk-go"
	"github.com/irisnet/core-sdk-go/common/crypto"
	"github.com/irisnet/core-sdk-go/types"
	sdktypes "github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/core-sdk-go/types/store"
)

func TestClient(t *testing.T) {

	nodeURI := "tcp://localhost:26657"
	grpcAddr := "localhost:9090"
	chainID := "wenchangchain"

	bech32AddressPrefix := sdktypes.AddrPrefixCfg{
		AccountAddr:   "iaa",
		ValidatorAddr: "iva",
		ConsensusAddr: "ica",
		AccountPub:    "iap",
		ValidatorPub:  "ivp",
		ConsensusPub:  "icp",
	}
	options := []sdktypes.Option{
		sdktypes.KeyDAOOption(store.NewMemory(nil)),
		sdktypes.TimeoutOption(10),
		sdktypes.KeyManagerOption(crypto.NewKeyManager()),
		sdktypes.Bech32AddressPrefixOption(&bech32AddressPrefix),
		sdktypes.BIP44PathOption(""),
	}
	cfg, err := types.NewClientConfig(nodeURI, grpcAddr, chainID, options...)
	if err != nil {
		panic(err)
	}

	sdkClient := sdk.NewClient(cfg)
	fmt.Println(sdkClient.EncodingConfig().TxConfig)

	//cli := NewClient(nil, txCfg)
	//txData := "0xf901470701831e8480941a6640c32b7e6413e839e9dfdb53970ee809b7fb80b8e4990711900000000000000000000000005892e7eeaea5ba624f5ba2900dbab8d2ea36d62b000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000001a687474703a2f2f6578616d706c652e636f6d2f746f6b656e2f33000000000000000000000000000000000000000000000000000000000000000000000000000974657374686173683200000000000000000000000000000000000000000000008209b2a05dc08aff9f0dac1ed240435510bdd53d8f8eb3f95c44a44f874e9e33ffd2407aa06090bf14d011822f2ff252081684c8e31a5e87d8735ee2f21acc32e87d28304f"
	//feePayer := "0x4579DB44FD3A6F645194058914E0A8D5E8F20DB8"
	//evmDenom := "ugas"
	//rawTx, err := cli.BuildEvmTx(txData, feePayer, evmDenom)
	//if err != nil {
	//	return
	//}
	//t.Log(rawTx)
}
