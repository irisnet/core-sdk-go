package gov

import (
	"context"
	"strconv"

	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/irisnet/core-sdk-go/codec"
	codetypes "github.com/irisnet/core-sdk-go/codec/types"
	"github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/core-sdk-go/types/errors"
	"github.com/irisnet/core-sdk-go/types/query"
)

type govClient struct {
	types.BaseClient
	codec.Codec
}

func NewClient(baseClient types.BaseClient, codec codec.Codec) Client {
	return govClient{
		BaseClient: baseClient,
		Codec:      codec,
	}
}

func (gc govClient) Name() string {
	return ModuleName
}

func (gc govClient) RegisterInterfaceTypes(registry codetypes.InterfaceRegistry) {
	RegisterInterfaces(registry)
}

func (gc govClient) SubmitProposal(request SubmitProposalRequest, baseTx types.BaseTx) (uint64, ctypes.ResultTx, error) {
	proposer, err := gc.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return 0, ctypes.ResultTx{}, errors.Wrap(ErrQueryAddress, err.Error())
	}

	deposit, err := gc.ToMinCoin(request.InitialDeposit...)
	if err != nil {
		return 0, ctypes.ResultTx{}, errors.Wrapf(ErrToMinCoin, err.Error())
	}

	content := ContentFromProposalType(request.Title, request.Description, request.Type)
	msg, err := NewMsgSubmitProposal(content, deposit, proposer)
	if err != nil {
		return 0, ctypes.ResultTx{}, errors.Wrap(ErrNewMsgSubmitProposal, err.Error())
	}

	result, err := gc.BuildAndSend([]types.Msg{msg}, baseTx)
	if err != nil {
		return 0, ctypes.ResultTx{}, errors.Wrap(ErrBuildAndSend, err.Error())
	}

	proposalIdStr, err := types.StringifyEvents(result.TxResult.Events).GetValue(types.EventTypeSubmitProposal, AttributeKeyProposalID)
	if err != nil {
		return 0, result, errors.Wrap(ErrEvensGetValue, err.Error())
	}

	proposalId, err := strconv.Atoi(proposalIdStr)
	if err != nil {
		return 0, result, errors.Wrap(ErrStrconvAtoi, err.Error())
	}
	return uint64(proposalId), result, err
}

func (gc govClient) Deposit(request DepositRequest, baseTx types.BaseTx) (ctypes.ResultTx, error) {
	depositor, err := gc.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return ctypes.ResultTx{}, errors.Wrap(ErrQueryAddress, err.Error())
	}

	amount, err := gc.ToMinCoin(request.Amount...)
	if err != nil {
		return ctypes.ResultTx{}, errors.Wrap(ErrToMinCoin, err.Error())
	}

	msg := &MsgDeposit{
		ProposalId: request.ProposalId,
		Depositor:  depositor.String(),
		Amount:     amount,
	}
	return gc.BuildAndSend([]types.Msg{msg}, baseTx)
}

// about VoteRequest.Option see  VoteOption_value
func (gc govClient) Vote(request VoteRequest, baseTx types.BaseTx) (ctypes.ResultTx, error) {
	voter, err := gc.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return ctypes.ResultTx{}, errors.Wrap(ErrQueryAddress, err.Error())
	}

	option := VoteOption_value[request.Option]
	msg := &MsgVote{
		ProposalId: request.ProposalId,
		Voter:      voter.String(),
		Option:     VoteOption(option),
	}
	return gc.BuildAndSend([]types.Msg{msg}, baseTx)
}

func (gc govClient) QueryProposal(proposalId uint64) (QueryProposalResp, error) {
	conn, err := gc.GenConn()

	if err != nil {
		return QueryProposalResp{}, errors.Wrap(ErrGenConn, err.Error())
	}

	res, err := NewQueryClient(conn).Proposal(
		context.Background(),
		&QueryProposalRequest{
			ProposalId: proposalId,
		})
	if err != nil {
		return QueryProposalResp{}, errors.Wrap(ErrQueryProposal, err.Error())
	}
	return res.Proposal.Convert().(QueryProposalResp), nil
}

// if proposalStatus is nil will return all status's proposals
// about proposalStatus see VoteOption_value
func (gc govClient) QueryProposals(proposalStatus string) ([]QueryProposalResp, error) {
	conn, err := gc.GenConn()

	if err != nil {
		return nil, errors.Wrap(ErrGenConn, err.Error())
	}

	res, err := NewQueryClient(conn).Proposals(
		context.Background(),
		&QueryProposalsRequest{
			ProposalStatus: ProposalStatus(VoteOption_value[proposalStatus]),
			Pagination: &query.PageRequest{
				Offset:     0,
				Limit:      100,
				CountTotal: true,
			},
		})
	if err != nil {
		return nil, errors.Wrap(ErrQueryProposal, err.Error())
	}
	return Proposals(res.Proposals).Convert().([]QueryProposalResp), nil
}

// about QueryVoteResp.Option see VoteOption_name
func (gc govClient) QueryVote(proposalId uint64, voter string) (QueryVoteResp, error) {
	conn, err := gc.GenConn()

	if err != nil {
		return QueryVoteResp{}, errors.Wrap(ErrGenConn, err.Error())
	}

	res, err := NewQueryClient(conn).Vote(
		context.Background(),
		&QueryVoteRequest{
			ProposalId: proposalId,
			Voter:      voter,
		})
	if err != nil {
		return QueryVoteResp{}, errors.Wrap(ErrQueryVote, err.Error())
	}
	return res.Vote.Convert().(QueryVoteResp), nil
}

func (gc govClient) QueryVotes(proposalId uint64) ([]QueryVoteResp, error) {
	conn, err := gc.GenConn()

	if err != nil {
		return nil, errors.Wrap(ErrGenConn, err.Error())
	}

	res, err := NewQueryClient(conn).Votes(
		context.Background(),
		&QueryVotesRequest{
			ProposalId: proposalId,
			Pagination: &query.PageRequest{
				Offset:     0,
				Limit:      100,
				CountTotal: true,
			},
		})
	if err != nil {
		return nil, errors.Wrap(ErrQueryVotes, err.Error())
	}
	return Votes(res.Votes).Convert().([]QueryVoteResp), nil
}

// QueryParams params_type("voting", "tallying", "deposit"), if don't pass will return all params_typ res
func (gc govClient) QueryParams(paramsType string) (QueryParamsResp, error) {
	conn, err := gc.GenConn()

	if err != nil {
		return QueryParamsResp{}, errors.Wrap(ErrGenConn, err.Error())
	}

	res, err := NewQueryClient(conn).Params(
		context.Background(),
		&QueryParamsRequest{
			ParamsType: paramsType,
		},
	)
	if err != nil {
		return QueryParamsResp{}, errors.Wrap(ErrQueryParams, err.Error())
	}
	return res.Convert().(QueryParamsResp), nil
}

func (gc govClient) QueryDeposit(proposalId uint64, depositor string) (QueryDepositResp, error) {
	conn, err := gc.GenConn()

	if err != nil {
		return QueryDepositResp{}, errors.Wrap(ErrGenConn, err.Error())
	}

	res, err := NewQueryClient(conn).Deposit(
		context.Background(),
		&QueryDepositRequest{
			ProposalId: proposalId,
			Depositor:  depositor,
		},
	)
	if err != nil {
		return QueryDepositResp{}, errors.Wrap(ErrQueryDeposit, err.Error())
	}
	return res.Deposit.Convert().(QueryDepositResp), nil
}

func (gc govClient) QueryDeposits(proposalId uint64) ([]QueryDepositResp, error) {
	conn, err := gc.GenConn()

	if err != nil {
		return nil, errors.Wrap(ErrGenConn, err.Error())
	}

	res, err := NewQueryClient(conn).Deposits(
		context.Background(),
		&QueryDepositsRequest{
			ProposalId: proposalId,
			Pagination: &query.PageRequest{
				Offset:     0,
				Limit:      100,
				CountTotal: true,
			},
		},
	)
	if err != nil {
		return nil, errors.Wrap(ErrQueryDeposit, err.Error())
	}
	return Deposits(res.Deposits).Convert().([]QueryDepositResp), nil
}

func (gc govClient) QueryTallyResult(proposalId uint64) (QueryTallyResultResp, error) {
	conn, err := gc.GenConn()

	if err != nil {
		return QueryTallyResultResp{}, errors.Wrap(ErrGenConn, err.Error())
	}

	res, err := NewQueryClient(conn).TallyResult(
		context.Background(),
		&QueryTallyResultRequest{
			ProposalId: proposalId,
		},
	)
	if err != nil {
		return QueryTallyResultResp{}, errors.Wrap(ErrQueryTallyResult, err.Error())
	}
	return res.Tally.Convert().(QueryTallyResultResp), nil
}
