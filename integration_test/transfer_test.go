package integration_test

import (
	"context"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/irisnet/core-sdk-go/modules/ibc/transfer"
	sdk "github.com/irisnet/core-sdk-go/types"
)

func (s IntegrationTestSuite) TestTransfer() {
	baseTx := sdk.BaseTx{
		From:     s.Account().Name,
		Gas:      200000,
		Memo:     "uiris",
		Mode:     sdk.Commit,
		Password: s.Account().Password,
	}
	sender, err2 := s.QueryAddress(baseTx.From, baseTx.Password)
	if err2 != nil {
		return
	}
	status, err := s.Status(context.Background())
	if err != nil {
		return
	}
	height := status.SyncInfo.LatestBlockHeight
	Request := transfer.TransferRequest{
		SourcePort:       "transfer",
		SourceChannel:    "channel-0",
		Token:            sdk.NewCoin("uiris", sdk.NewInt(1024)),
		Sender:           sender.String(),
		Receiver:         "iaa10njupdhmnyma2s7ghcapgtnw9kzg9gkjdylyla",
		TimeoutHeight:    transfer.Height{RevisionHeight: uint64(height + 128)},
		TimeoutTimestamp: uint64(time.Now().Add(time.Minute * 8).UnixNano()),
	}

	result, err := s.Transfer.Transfer(Request, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), result.Hash)
}
