package integration_test

import (
	"github.com/irisnet/core-sdk-go/feegrant"
	"github.com/irisnet/core-sdk-go/types"
	"github.com/stretchr/testify/require"
	"time"
)

func (s IntegrationTestSuite) TestFeeGrant() {
	cases := []SubTest{
		{
			"TestGrant",
			grant,
		},
	}

	for _, t := range cases {
		s.Run(t.testName, func() { t.testCase(s) })
	}
}

func grant(s IntegrationTestSuite) {
	to, _ := types.AccAddressFromBech32("iaa1pn9dv6hh5lhy4lpya3scl95r3fhrdfv26kjfgh")
	atom := types.NewCoins(types.NewInt64Coin("uirita", 555))
	threeHours := time.Now().Add(3 * time.Hour)
	basic := &feegrant.BasicAllowance{
		SpendLimit: atom,
		Expiration: &threeHours,
	}

	baseTx := types.BaseTx{
		From:               s.Account().Name,
		Gas:                200000,
		Memo:               "TEST",
		Mode:               types.Commit,
		Password:           s.Account().Password,
		SimulateAndExecute: false,
		GasAdjustment:      1.5,
	}

	result, err := s.FeeGrant.GrantAllowance(s.Account().Address, to, basic, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), result.Hash)
}
