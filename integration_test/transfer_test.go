package integration_test

// import (
// 	"github.com/stretchr/testify/require"

// 	"github.com/irisnet/core-sdk-go/modules/ibc/transfer"
// 	sdk "github.com/irisnet/core-sdk-go/types"
// )

// func (s IntegrationTestSuite) TestTransfer() {
// 	baseTx := sdk.BaseTx{
// 		From:     s.Account().Name,
// 		Gas:      200000,
// 		Memo:     "test",
// 		Mode:     sdk.Commit,
// 		Password: s.Account().Password,
// 	}
// 	Request := transfer.TransferRequest{}

// 	result, err := s.Transfer.Transfer(Request, baseTx)
// 	require.NoError(s.T(), err)
// 	require.NotEmpty(s.T(), result.Hash)

// }
