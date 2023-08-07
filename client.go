package sdk

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/irisnet/core-sdk-go/client"
	bank2 "github.com/irisnet/core-sdk-go/modules/bank"
	"github.com/irisnet/core-sdk-go/modules/feegrant"
	sdk "github.com/irisnet/core-sdk-go/types"
	"github.com/tendermint/tendermint/libs/log"
)

type Client struct {
	logger         log.Logger
	moduleManager  map[string]sdk.Module
	encodingConfig sdk.EncodingConfig
	sdk.BaseClient
	Bank     bank2.Client
	FeeGrant feegrant.Client
}

func NewClient(cfg sdk.ClientConfig) Client {
	encodingConfig := sdk.MakeEncodingConfig()

	// create a instance of baseClient
	baseClient := client.NewBaseClient(cfg, encodingConfig)
	bankClient := bank2.NewClient(baseClient, encodingConfig.Marshaler)
	feeGrantClient := feegrant.NewClient(baseClient, encodingConfig.Marshaler)

	client := Client{
		BaseClient:     baseClient,
		moduleManager:  make(map[string]sdk.Module),
		encodingConfig: encodingConfig,
		Bank:           bankClient,
		FeeGrant:       feeGrantClient,
	}
	client.RegisterModule(
		bankClient,
		feeGrantClient,
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
	return client.encodingConfig.Marshaler
}

func (client *Client) EncodingConfig() sdk.EncodingConfig {
	return client.encodingConfig
}

func (client *Client) Manager() sdk.BaseClient {
	return client.BaseClient
}

func (client *Client) RegisterModule(ms ...sdk.Module) {
	for _, m := range ms {
		m.RegisterInterfaceTypes(client.encodingConfig.InterfaceRegistry)
	}
}

func (client *Client) Module(name string) sdk.Module {
	return client.moduleManager[name]
}
