package gov

import (
	"github.com/irisnet/core-sdk-go/types/errors"
)

const Codespace = ModuleName

var (
	ErrDuplicatedVoteOption      = errors.Register(Codespace, 1, "duplicated vote option")
	ErrInvalidTotalWeight        = errors.Register(Codespace, 2, "invalid total weight")
	ErrInvalidWeightedVoteOption = errors.Register(Codespace, 3, "invalid weighted vote option")
	ErrInvalidVoteOption         = errors.Register(Codespace, 4, "invalid vote option")
	ErrQueryAddress              = errors.Register(Codespace, 5, "query address error")
	ErrNewMsgSubmitProposal      = errors.Register(Codespace, 6, "NewMsgSubmitProposal error")
	ErrBuildAndSend              = errors.Register(Codespace, 7, "BuildAndSend error")
	ErrEvensGetValue             = errors.Register(Codespace, 8, "EvensGetValue error")
	ErrStrconvAtoi               = errors.Register(Codespace, 9, "StrconvAtoi error")
	ErrQueryProposal             = errors.Register(Codespace, 10, "QueryProposal error")
	ErrQueryVote                 = errors.Register(Codespace, 11, "QueryVote error")
	ErrQueryVotes                = errors.Register(Codespace, 12, "QueryVotes error")
	ErrQueryDeposit              = errors.Register(Codespace, 13, "QueryDeposit error")
	ErrQueryTallyResult          = errors.Register(Codespace, 14, "QueryTallyResult error")
	ErrInvalidTitle              = errors.Register(Codespace, 15, "invalid title")
	ErrInvalidProposer           = errors.Register(Codespace, 16, "invalid proposer")
	ErrInvalidCoin               = errors.Register(Codespace, 17, "invalid coin")
	ErrMissingContent            = errors.Register(Codespace, 18, "MissingContent error")
	ErrMissingProposer           = errors.Register(Codespace, 19, "MissingProposer error")
	ErrInvalidAmount             = errors.Register(Codespace, 20, "invalid amount")
	ErrInvalidDescription        = errors.Register(Codespace, 21, "invalid description")
	ErrToMinCoin                 = errors.Register(Codespace, 22, "ToMinCoin error")
	ErrQueryParams               = errors.Register(Codespace, 23, "query params  error")
	ErrGenConn                   = errors.Register(Codespace, 24, "generate conn error")
)
