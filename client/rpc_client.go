package client

import (
	"context"
	"fmt"

	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/libs/log"
	rpc "github.com/tendermint/tendermint/rpc/client"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/core-sdk-go/types/errors"
	"github.com/irisnet/core-sdk-go/uuid"
)

type rpcClient struct {
	rpc.Client
	log.Logger
	txDecoder types.TxDecoder
}

func NewRPCClient(
	remote string,
	txDecoder types.TxDecoder,
	logger log.Logger,
	timeout uint,
) types.TmClient {
	client, err := rpchttp.NewWithTimeout(remote, "/websocket", timeout)
	if err != nil {
		panic(err)
	}

	_ = client.Start()
	return rpcClient{
		Client:    client,
		Logger:    logger,
		txDecoder: txDecoder,
	}
}

// =============================================================================
// SubscribeNewBlock implement WSClient interface
func (r rpcClient) SubscribeNewBlock(builder *types.EventQueryBuilder, handler types.EventNewBlockHandler) (types.Subscription, error) {
	if builder == nil {
		builder = types.NewEventQueryBuilder()
	}
	builder.AddCondition(types.Cond(types.TypeKey).EQ(tmtypes.EventNewBlock))
	query := builder.Build()

	return r.SubscribeAny(query, func(data types.EventData) {
		handler(data.(types.EventDataNewBlock))
	})
}

// SubscribeTx implement WSClient interface
func (r rpcClient) SubscribeTx(builder *types.EventQueryBuilder, handler types.EventTxHandler) (types.Subscription, error) {
	if builder == nil {
		builder = types.NewEventQueryBuilder()
	}
	query := builder.AddCondition(types.Cond(types.TypeKey).EQ(types.TxValue)).Build()
	return r.SubscribeAny(query, func(data types.EventData) {
		handler(data.(types.EventDataTx))
	})
}

func (r rpcClient) SubscribeNewBlockHeader(handler types.EventNewBlockHeaderHandler) (types.Subscription, error) {
	query := tmtypes.QueryForEvent(tmtypes.EventNewBlockHeader).String()
	return r.SubscribeAny(query, func(data types.EventData) {
		handler(data.(types.EventDataNewBlockHeader))
	})
}

func (r rpcClient) SubscribeValidatorSetUpdates(handler types.EventValidatorSetUpdatesHandler) (types.Subscription, error) {
	query := tmtypes.QueryForEvent(tmtypes.EventValidatorSetUpdates).String()
	return r.SubscribeAny(query, func(data types.EventData) {
		handler(data.(types.EventDataValidatorSetUpdates))
	})
}

func (r rpcClient) Resubscribe(subscription types.Subscription, handler types.EventHandler) (err error) {
	_, err = r.SubscribeAny(subscription.Query, handler)
	return
}

func (r rpcClient) Unsubscribe(subscription types.Subscription) error {
	r.Info("end to subscribe event", "query", subscription.Query, "subscriber", subscription.ID)
	err := r.Client.Unsubscribe(subscription.Ctx, subscription.ID, subscription.Query)
	if err != nil {
		r.Error("unsubscribe failed", "query", subscription.Query, "subscriber", subscription.ID, "errMsg", err.Error())
		return errors.Wrap(errors.ErrTodo, err.Error())
	}
	return nil
}

func (r rpcClient) SubscribeAny(query string, handler types.EventHandler) (subscription types.Subscription, err error) {
	ctx := context.Background()
	subscriber := getSubscriber()
	ch, err := r.Subscribe(ctx, subscriber, query, 0)
	if err != nil {
		return subscription, errors.Wrap(errors.ErrTodo, err.Error())
	}

	r.Info("subscribe event", "query", query, "subscriber", subscription.ID)

	subscription = types.Subscription{
		Ctx:   ctx,
		Query: query,
		ID:    subscriber,
	}
	go func() {
		for {
			data := <-ch
			go func() {
				defer errors.CatchPanic(func(errMsg string) {
					r.Error("unsubscribe failed", "query", subscription.Query, "subscriber", subscription.ID, "errMsg", err.Error())
				})

				switch data := data.Data.(type) {
				case tmtypes.EventDataTx:
					handler(r.parseTx(data))
					return
				case tmtypes.EventDataNewBlock:
					handler(r.parseNewBlock(data))
					return
				case tmtypes.EventDataNewBlockHeader:
					handler(r.parseNewBlockHeader(data))
					return
				case tmtypes.EventDataValidatorSetUpdates:
					handler(r.parseValidatorSetUpdates(data))
					return
				default:
					handler(data)
				}
			}()
		}
	}()
	return
}

func (r rpcClient) parseTx(data types.EventData) types.EventDataTx {
	dataTx := data.(tmtypes.EventDataTx)
	tx, err := r.txDecoder(dataTx.Tx)
	if err != nil {
		return types.EventDataTx{}
	}

	hash := types.HexBytes(tmhash.Sum(dataTx.Tx)).String()
	result := types.TxResult{
		Code:      dataTx.Result.Code,
		Log:       dataTx.Result.Log,
		GasWanted: dataTx.Result.GasWanted,
		GasUsed:   dataTx.Result.GasUsed,
		Events:    types.StringifyEvents(dataTx.Result.Events),
	}

	return types.EventDataTx{
		Hash:   hash,
		Height: dataTx.Height,
		Index:  dataTx.Index,
		Tx:     tx,
		Result: result,
	}
}

func (r rpcClient) parseNewBlock(data types.EventData) types.EventDataNewBlock {
	block := data.(tmtypes.EventDataNewBlock)
	return types.EventDataNewBlock{
		Block: types.ParseBlock(r.txDecoder, block.Block),
		ResultBeginBlock: types.ResultBeginBlock{
			Events: types.StringifyEvents(block.ResultBeginBlock.Events),
		},
		ResultEndBlock: types.ResultEndBlock{
			Events:           types.StringifyEvents(block.ResultEndBlock.Events),
			ValidatorUpdates: types.ParseValidatorUpdate(block.ResultEndBlock.ValidatorUpdates),
		},
	}
}

func (r rpcClient) parseNewBlockHeader(data types.EventData) types.EventDataNewBlockHeader {
	blockHeader := data.(tmtypes.EventDataNewBlockHeader)
	return types.EventDataNewBlockHeader{
		Header: blockHeader.Header,
		ResultBeginBlock: types.ResultBeginBlock{
			Events: types.StringifyEvents(blockHeader.ResultBeginBlock.Events),
		},
		ResultEndBlock: types.ResultEndBlock{
			Events:           types.StringifyEvents(blockHeader.ResultEndBlock.Events),
			ValidatorUpdates: types.ParseValidatorUpdate(blockHeader.ResultEndBlock.ValidatorUpdates),
		},
	}
}

func (r rpcClient) parseValidatorSetUpdates(data types.EventData) types.EventDataValidatorSetUpdates {
	validatorSet := data.(tmtypes.EventDataValidatorSetUpdates)
	return types.EventDataValidatorSetUpdates{
		ValidatorUpdates: types.ParseValidators(validatorSet.ValidatorUpdates),
	}
}

func getSubscriber() string {
	subscriber := "core-sdk-go"
	id, err := uuid.NewV1()
	if err == nil {
		subscriber = fmt.Sprintf("%s-%s", subscriber, id.String())
	}
	return subscriber
}
