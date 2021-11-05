package sdk

import (
	"github.com/tendermint/tendermint/libs/log"

	"github.com/irisnet/core-sdk-go/client"
	keys "github.com/irisnet/core-sdk-go/client"
	"github.com/irisnet/core-sdk-go/codec"
	cryptotypes "github.com/irisnet/core-sdk-go/codec/types"
	cryptocodec "github.com/irisnet/core-sdk-go/crypto/codec"
	"github.com/irisnet/core-sdk-go/modules/bank"
	"github.com/irisnet/core-sdk-go/modules/gov"
	"github.com/irisnet/core-sdk-go/modules/ibc/transfer"
	"github.com/irisnet/core-sdk-go/modules/staking"
	"github.com/irisnet/core-sdk-go/types"
	txtypes "github.com/irisnet/core-sdk-go/types/tx"
)

type Client struct {
	types.BaseClient
	logger         log.Logger
	moduleManager  map[string]types.Module
	encodingConfig types.EncodingConfig
	Bank           bank.Client
	Key            keys.Client
	Staking        staking.Client
	Gov            gov.Client
	Transfer       transfer.Client
}

func NewClient(cfg types.ClientConfig) Client {
	encodingConfig := makeEncodingConfig()

	// create a instance of baseClient
	baseClient := client.NewBaseClient(cfg, encodingConfig, nil)
	bankClient := bank.NewClient(baseClient, encodingConfig.Codec)
	keysClient := keys.NewKeysClient(cfg, baseClient)
	transferClient := transfer.NewClient(baseClient, encodingConfig.Codec)
	stakingClient := staking.NewClient(baseClient, encodingConfig.Codec)
	govClient := gov.NewClient(baseClient, encodingConfig.Codec)

	client := Client{
		logger:         baseClient.Logger(),
		BaseClient:     baseClient,
		moduleManager:  make(map[string]types.Module),
		encodingConfig: encodingConfig,
		Bank:           bankClient,
		Key:            keysClient,
		Staking:        stakingClient,
		Gov:            govClient,
		Transfer:       transferClient,
	}
	client.RegisterModule(
		bankClient,
		stakingClient,
		govClient,
		transferClient,
	)
	return client
}

func (client *Client) SetLogger(logger log.Logger) {
	client.BaseClient.SetLogger(logger)
}

func (client *Client) Codec() *codec.LegacyAmino {
	return client.encodingConfig.Amino
}

func (client *Client) AppCodec() codec.Codec {
	return client.encodingConfig.Codec
}

func (client *Client) EncodingConfig() types.EncodingConfig {
	return client.encodingConfig
}

func (client *Client) Manager() types.BaseClient {
	return client.BaseClient
}

func (client *Client) RegisterModule(ms ...types.Module) {
	for _, m := range ms {
		m.RegisterInterfaceTypes(client.encodingConfig.InterfaceRegistry)
	}
}

func (client *Client) Module(name string) types.Module {
	return client.moduleManager[name]
}

func makeEncodingConfig() types.EncodingConfig {
	amino := codec.NewLegacyAmino()
	interfaceRegistry := cryptotypes.NewInterfaceRegistry()
	codec := codec.NewProtoCodec(interfaceRegistry)
	txCfg := txtypes.NewTxConfig(codec, txtypes.DefaultSignModes)

	encodingConfig := types.EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Codec:             codec,
		TxConfig:          txCfg,
		Amino:             amino,
	}
	RegisterLegacyAminoCodec(encodingConfig.Amino)
	RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}

// RegisterLegacyAminoCodec registers the sdk message type.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*types.Msg)(nil), nil)
	cdc.RegisterInterface((*types.Tx)(nil), nil)
	cryptocodec.RegisterCrypto(cdc)
}

// RegisterInterfaces registers the sdk message type.
func RegisterInterfaces(registry cryptotypes.InterfaceRegistry) {
	registry.RegisterInterface("cosmos.v1beta1.Msg", (*types.Msg)(nil))
	txtypes.RegisterInterfaces(registry)
	cryptocodec.RegisterInterfaces(registry)
}
