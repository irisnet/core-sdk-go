package gov

import (
	"github.com/irisnet/core-sdk-go/types/errors"
)

const Codespace = ModuleName

var (
	ErrDuplicatedVoteOption      = errors.Register(Codespace, 2, "duplicated vote option")
	ErrInvalidTotalWeight        = errors.Register(Codespace, 2, "invalid total weight")
	ErrInvalidWeightedVoteOption = errors.Register(Codespace, 2, "invalid weighted vote option")
	ErrInvalidVoteOption         = errors.Register(Codespace, 5, "invalid vote option")
	ErrQueryAddress              = errors.Register(Codespace, 5, "query address error")
	ErrNewMsgSubmitProposal      = errors.Register(Codespace, 5, "NewMsgSubmitProposal error")
	ErrBuildAndSend              = errors.Register(Codespace, 5, "BuildAndSend error")
	ErrEvensGetValue             = errors.Register(Codespace, 5, "EvensGetValue error")
	ErrStrconvAtoi               = errors.Register(Codespace, 5, "StrconvAtoi error")
	ErrQueryProposal             = errors.Register(Codespace, 5, "QueryProposal error")
	ErrQueryVote                 = errors.Register(Codespace, 5, "QueryVote error")
	ErrQueryVotes                = errors.Register(Codespace, 5, "QueryVotes error")
	ErrQueryDeposit              = errors.Register(Codespace, 5, "QueryDeposit error")
	ErrQueryTallyResult          = errors.Register(Codespace, 5, "QueryTallyResult error")
	ErrInvalidTitle              = errors.Register(Codespace, 5, "invalid title")
	ErrInvalidProposer           = errors.Register(Codespace, 5, "invalid proposer")
	ErrInvalidCoin               = errors.Register(Codespace, 5, "invalid coin")
	ErrMissingContent            = errors.Register(Codespace, 5, "MissingContent error")
	ErrMissingProposer           = errors.Register(Codespace, 5, "MissingProposer error")
	ErrInvalidAmount             = errors.Register(Codespace, 1, "invalid amount")
	ErrInvalidDescription        = errors.Register(Codespace, 4, "invalid description")
	ErrToMinCoin                 = errors.Register(Codespace, 6, "ToMinCoin error")
	ErrQueryParams               = errors.Register(Codespace, 21, "query params  error")
	ErrGenConn                   = errors.Register(Codespace, 22, "generate conn error")
)
