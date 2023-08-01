package types

import (
	"errors"

	"github.com/cosmos/cosmos-sdk/types/query"
)

const (
	CountTotalErrMsg = "pageRequest error: CountTotal is not supported, must be false"
	LimitErrMsg      = "pageRequest error: Limit cannot be empty and cannot exceed 100"
	OffsetKeyErrMsg  = "pageRequest error: Only one Offset or Key is allowed"
)

func FormatPageRequest(pageReq *query.PageRequest) (*query.PageRequest, Error) {
	if pageReq == nil {
		return &query.PageRequest{
			Offset:     0,
			Limit:      100,
			CountTotal: false,
		}, nil
	}

	if pageReq.CountTotal {
		return pageReq, Wrap(errors.New(CountTotalErrMsg))
	}

	if pageReq.Limit == 0 || pageReq.Limit > 100 {
		return pageReq, Wrap(errors.New(LimitErrMsg))
	}

	if pageReq.Offset != 0 && len(pageReq.Key) > 0 {
		return pageReq, Wrap(errors.New(OffsetKeyErrMsg))
	}

	return pageReq, nil
}
