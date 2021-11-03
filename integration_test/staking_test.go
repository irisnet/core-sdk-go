package integration_test

import (
	"context"

	tmed25519 "github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/irisnet/core-sdk-go/crypto/codec"

	"github.com/stretchr/testify/require"

	"github.com/irisnet/core-sdk-go/modules/staking"
	sdk "github.com/irisnet/core-sdk-go/types"
)

func (s IntegrationTestSuite) TestStaking() {
	cases := []SubTest{
		{"TestCreateAndEdit", testCreateAndEdit},
		{"TestStaking", testStaking},
		{"TestQueryHistoricalInfo", queryHistoricalInfo},
		{"TestQueryPool", queryPool},
		{"TestQueryParams", queryParams},
	}

	for _, t := range cases {
		s.Run(t.testName, func() {
			t.testCase(s)
		})
	}
}

// this need another node to test
func testCreateAndEdit(s IntegrationTestSuite) {
	// send createValidator tx
	baseTx := sdk.BaseTx{
		From:     s.Account().Name,
		Gas:      200000,
		Memo:     "test",
		Mode:     sdk.Commit,
		Password: s.Account().Password,
	}

	name, password := s.RandStringOfLength(8), s.RandStringOfLength(8)
	address, mnemonic, err := s.Key.Add(name, password)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), address)
	require.NotEmpty(s.T(), mnemonic)
	queryAddress, err := s.QueryAddress(name, password)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), queryAddress)

	_, err = s.Bank.Send(address, sdk.NewDecCoins(sdk.NewDecCoin("iris", sdk.NewInt(10))), baseTx)
	require.NoError(s.T(), err)

	privkey := tmed25519.GenPrivKey()
	public255 := privkey.PubKey()
	sdkPK, err := codec.FromTmPubKeyInterface(public255)
	require.NoError(s.T(), err)
	bz, err := s.AppCodec().MarshalInterfaceJSON(sdkPK)
	require.NoError(s.T(), err)

	stakingBaseTx := sdk.BaseTx{
		From:     name,
		Gas:      200000,
		Memo:     "test",
		Mode:     sdk.Commit,
		Password: password,
	}

	rate := sdk.MustNewDecFromStr("0.1")
	maxRate := sdk.MustNewDecFromStr("0.1")
	maxChangeRate := sdk.MustNewDecFromStr("0.01")
	minSelfDelegation := sdk.OneInt()
	value, _ := sdk.ParseDecCoin("10uiris")
	req1 := staking.CreateValidatorRequest{
		Moniker:           "haha",
		Rate:              rate,
		MaxRate:           maxRate,
		MaxChangeRate:     maxChangeRate,
		MinSelfDelegation: minSelfDelegation,
		Pubkey:            string(bz),
		Value:             value,
	}
	res, err := s.Staking.CreateValidator(req1, stakingBaseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), res.Hash)

	// send editValidator tx
	commissionRate := sdk.MustNewDecFromStr("0.1")
	minSelfDelegation = sdk.NewInt(2)
	req2 := staking.EditValidatorRequest{
		Moniker:           "haha",
		Identity:          "identity",
		Website:           "website",
		SecurityContact:   "abbccdd",
		Details:           "fadsfas",
		CommissionRate:    commissionRate,
		MinSelfDelegation: minSelfDelegation,
	}
	res, err = s.Staking.EditValidator(req2, stakingBaseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), res.Hash)
}

func testStaking(s IntegrationTestSuite) {
	// ================================ about delegate ==============================
	delegateAddr := s.Account().Address.String()
	baseTx := sdk.BaseTx{
		From:     s.Account().Name,
		Gas:      200000,
		Memo:     "test",
		Mode:     sdk.Commit,
		Password: s.Account().Password,
	}

	// queries all validators that match the given status.
	validatorsResp, err := s.Staking.QueryValidators("", 1, 100)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), validatorsResp.Validators)

	// queries validator info for given validator address.
	validatorAddr := validatorsResp.Validators[0].OperatorAddress
	validatorResp, err := s.Staking.QueryValidator(validatorAddr)
	require.NoError(s.T(), err)
	require.Equal(s.T(), validatorAddr, validatorResp.OperatorAddress)

	// send Delegate tx
	amount, _ := sdk.ParseDecCoin("10iris")
	delegateReq := staking.DelegateRequest{
		ValidatorAddr: validatorAddr,
		Amount:        amount,
	}
	res, err := s.Staking.Delegate(delegateReq, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), res.Hash)

	// queries delegate info for given validator delegator pair.
	delegation, err := s.Staking.QueryDelegation(delegateAddr, validatorAddr)
	require.NoError(s.T(), err)
	require.Equal(s.T(), delegateAddr, delegation.Delegation.DelegatorAddress)
	require.Equal(s.T(), validatorAddr, delegation.Delegation.ValidatorAddress)

	// queries delegate info for given validator
	delegationsToResp, err := s.Staking.QueryValidatorDelegations(validatorAddr, 1, 10)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), delegationsToResp.DelegationResponses)
	require.Greater(s.T(), delegationsToResp.Total, uint64(0))
	var exists bool
	for _, d := range delegationsToResp.DelegationResponses {
		if d.Delegation.DelegatorAddress == delegateAddr {
			exists = true
		}
	}
	require.True(s.T(), exists)

	// queries all delegations of a given delegator address.
	delegatorDelegations, err := s.Staking.QueryDelegatorDelegations(delegateAddr, 1, 10)
	require.NoError(s.T(), err)
	exists = false // init exists
	for _, d := range delegatorDelegations.DelegationResponses {
		if d.Delegation.ValidatorAddress == validatorAddr && d.Delegation.DelegatorAddress == delegateAddr {
			exists = true
		}
	}
	require.True(s.T(), exists)

	// queries all validators info for given delegator
	delegatorValidators, err := s.Staking.QueryDelegatorValidators(delegateAddr, 1, 10)
	require.NoError(s.T(), err)
	exists = false // init exists
	for _, v := range delegatorValidators.Validator {
		if v.OperatorAddress == validatorAddr {
			exists = true
		}
	}
	require.True(s.T(), exists)

	// queries validator info for given delegator validator pair.
	delegatorValidator, err := s.Staking.QueryDelegatorValidator(delegateAddr, validatorAddr)
	require.NoError(s.T(), err)
	require.Equal(s.T(), validatorAddr, delegatorValidator.OperatorAddress)

	// ================================ about unbonding ==============================
	// send Undelegate tx
	amount, _ = sdk.ParseDecCoin("10iris")
	undelegateReq := staking.UndelegateRequest{
		ValidatorAddr: validatorAddr,
		Amount:        amount,
	}
	res, err = s.Staking.Undelegate(undelegateReq, baseTx)
	require.NoError(s.T(), err)
	require.Greater(s.T(), res.Height, int64(1))

	// queries unbonding delegations of a validator.
	unbondingDelegations, err := s.Staking.QueryValidatorUnbondingDelegations(validatorAddr, 1, 10)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), unbondingDelegations.UnbondingResponses)
	exists = false // init exists
	for _, u := range unbondingDelegations.UnbondingResponses {
		if u.DelegatorAddress == delegateAddr && u.ValidatorAddress == validatorAddr {
			exists = true
		}
		require.NotEmpty(s.T(), u.Entries)
	}
	require.True(s.T(), exists)

	// queries unbonding info for given validator delegator pair.
	unbondingDelegation, err := s.Staking.QueryUnbondingDelegation(delegateAddr, validatorAddr)
	require.NoError(s.T(), err)
	require.Equal(s.T(), validatorAddr, unbondingDelegation.ValidatorAddress)
	require.Equal(s.T(), delegateAddr, unbondingDelegation.DelegatorAddress)

	// queries all unbonding delegations of a given delegator address.
	delegatorUnbondingDelegations, err := s.Staking.QueryDelegatorUnbondingDelegations(delegateAddr, 1, 10)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), res.Hash)
	exists = false // init exists
	for _, d := range delegatorUnbondingDelegations.UnbondingDelegations {
		if d.DelegatorAddress == delegateAddr && d.ValidatorAddress == validatorAddr {
			exists = true
		}
		require.NotEmpty(s.T(), d.Entries)
	}
	require.True(s.T(), exists)

	// ================================ about redelegate ==============================
	// send redelegate tx
	//amount, _ = sdk.ParseDecCoin("3000iris")
	// you can use another node to create a validator, then assgin newValidatorAddr in ValidatorDstAddress to send this txDelegatorAddress: delegatorAddr.String(),
	//newValidatorAddr := validatorAddr
	//redelegateReq := staking.BeginRedelegateRequest{
	//	ValidatorSrcAddress: validatorAddr,
	//	ValidatorDstAddress: newValidatorAddr,
	//	Amount:              amount,
	//}
	//res, err = s.Staking.BeginRedelegate(redelegateReq, baseTx)
	//require.NoError(s.T(), err)
	//require.NotEmpty(s.T(), res.Hash)
	//
	//// queries redelegations of given address.
	//redelegationsReq := staking.QueryRedelegationsReq{
	//	DelegatorAddr:    delegateAddr,
	//	SrcValidatorAddr: validatorAddr,
	//	DstValidatorAddr: newValidatorAddr,
	//	Page:             1,
	//	Size:             10,
	//}
	//redelegations, err := s.Staking.QueryRedelegations(redelegationsReq)
	//require.NoError(s.T(), err)
	//exists = false // init exists
	//for _, r := range redelegations.RedelegationResponses {
	//	if r.Redelegation.ValidatorSrcAddress == validatorAddr && r.Redelegation.ValidatorDstAddress == newValidatorAddr {
	//		exists = true
	//	}
	//	require.NotEmpty(s.T(), r.Entries)
	//}
	//require.True(s.T(), exists)
}

func queryHistoricalInfo(s IntegrationTestSuite) {
	// get latestBlockHeight at first
	status, err := s.Status(context.Background())
	require.NoError(s.T(), err)
	height := status.SyncInfo.LatestBlockHeight
	height -= 10

	res, err := s.Staking.QueryHistoricalInfo(height)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), res.Valset)

	valAddr := s.curValAddr()
	var flag bool
	for _, validator := range res.Valset {
		for _, valAdd := range valAddr {
			if validator.OperatorAddress == valAdd.OperatorAddress {
				flag = true
			}
		}
	}
	require.True(s.T(), flag)
}

func queryPool(s IntegrationTestSuite) {
	res, err := s.Staking.QueryPool()
	require.NoError(s.T(), err)
	require.Greater(s.T(), res.BondedTokens.Int64(), int64(0))
	require.GreaterOrEqual(s.T(), res.NotBondedTokens.Int64(), int64(0)) // NotBondedTokens can be 0
}

func queryParams(s IntegrationTestSuite) {
	// this params is irishub default params
	const (
		bondDenom         = "uiris"
		defaultHistorical = uint32(10000)
		MaxValidators     = uint32(100)
		MaxEntries        = uint32(7)
	)

	res, err := s.Staking.QueryParams()
	require.NoError(s.T(), err)
	require.Equal(s.T(), bondDenom, res.BondDenom)
	require.Equal(s.T(), defaultHistorical, res.HistoricalEntries)
	require.Equal(s.T(), MaxValidators, res.MaxValidators)
	require.Equal(s.T(), MaxEntries, res.MaxEntries)
}

func (s IntegrationTestSuite) curValAddr() []staking.QueryValidatorResp {
	// queries all validators that match the given status.
	validatorsResp, err := s.Staking.QueryValidators("", 1, 200)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), validatorsResp.Validators)
	return validatorsResp.Validators
}
