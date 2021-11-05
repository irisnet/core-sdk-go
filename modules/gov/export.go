package gov

import (
	"time"

	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/irisnet/core-sdk-go/types"
)

// expose Gov module api for user
type Client interface {
	types.Module
	SubmitProposal(request SubmitProposalRequest, baseTx types.BaseTx) (uint64, ctypes.ResultTx, error)
	Deposit(request DepositRequest, baseTx types.BaseTx) (ctypes.ResultTx, error)
	Vote(request VoteRequest, baseTx types.BaseTx) (ctypes.ResultTx, error)

	QueryProposal(proposalId uint64) (QueryProposalResp, error)
	QueryProposals(proposalStatus string) ([]QueryProposalResp, error)
	QueryVote(proposalId uint64, voter string) (QueryVoteResp, error)
	QueryVotes(proposalId uint64) ([]QueryVoteResp, error)
	QueryParams(paramsType string) (QueryParamsResp, error)
	QueryDeposit(proposalId uint64, depositor string) (QueryDepositResp, error)
	QueryDeposits(proposalId uint64) ([]QueryDepositResp, error)
	QueryTallyResult(proposalId uint64) (QueryTallyResultResp, error)
}

type SubmitProposalRequest struct {
	Title          string         `json:"title"`
	Description    string         `json:"description"`
	Type           string         `json:"type"`
	InitialDeposit types.DecCoins `json:"initial_deposit"`
}

type DepositRequest struct {
	ProposalId uint64         `json:"proposal_id"`
	Amount     types.DecCoins `json:"amount"`
}

type VoteRequest struct {
	ProposalId uint64 `json:"proposal_id"`
	Option     string `json:"option"`
}

type QueryProposalResp struct {
	ProposalId       uint64               `json:"proposal_id"`
	Content          Content              `json:"content"`
	Status           string               `json:"status"`
	FinalTallyResult QueryTallyResultResp `json:"final_tally_result"`
	SubmitTime       time.Time            `json:"submit_time"`
	DepositEndTime   time.Time            `json:"deposit_end_time"`
	TotalDeposit     types.Coins          `json:"total_deposit"`
	VotingStartTime  time.Time            `json:"voting_start_time"`
	VotingEndTime    time.Time            `json:"voting_end_time"`
}

type QueryVoteResp struct {
	ProposalId uint64 `json:"proposal_id"`
	Voter      string `json:"voter"`
	Option     int32  `json:"option"`
}

type (
	votingParams struct {
		VotingPeriod time.Duration `json:"voting_period"`
	}
	depositParams struct {
		MinDeposit       types.Coins   `json:"min_deposit"`
		MaxDepositPeriod time.Duration `json:"max_deposit_period"`
	}
	tallyParams struct {
		Quorum        types.Dec `json:"quorum"`
		Threshold     types.Dec `json:"threshold"`
		VetoThreshold types.Dec `json:"veto_threshold"`
	}
	QueryParamsResp struct {
		VotingParams  votingParams  `json:"voting_params"`
		DepositParams depositParams `json:"deposit_params"`
		TallyParams   tallyParams   `json:"tally_params"`
	}
)

type QueryDepositResp struct {
	ProposalId uint64      `json:"proposal_id"`
	Depositor  string      `json:"depositor"`
	Amount     types.Coins `json:"amount"`
}

type QueryTallyResultResp struct {
	Yes        types.Int `json:"yes"`
	Abstain    types.Int `json:"abstain"`
	No         types.Int `json:"no"`
	NoWithVeto types.Int `json:"no_with_veto"`
}
