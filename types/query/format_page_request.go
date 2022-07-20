package query

import sdk "github.com/irisnet/core-sdk-go/types"

const (
	CountTotalErrMsg = "pageRequest error: CountTotal is not supported, must be false"
	LimitErrMsg      = "pageRequest error: Limit cannot be empty and cannot exceed 100"
	OffsetKeyErrMsg  = "pageRequest error: Only one Offset or Key is allowed"
)

func FormatPageRequest(pageReq *PageRequest) (*PageRequest, sdk.Error) {
	if pageReq == nil {
		return &PageRequest{
			Offset:     0,
			Limit:      100,
			CountTotal: false,
		}, nil
	}

	if pageReq.CountTotal {
		return pageReq, sdk.Wrapf(CountTotalErrMsg)
	}

	if pageReq.Limit == 0 || pageReq.Limit > 100 {
		return pageReq, sdk.Wrapf(LimitErrMsg)
	}

	if pageReq.Offset != 0 && len(pageReq.Key) > 0 {
		return pageReq, sdk.Wrapf(OffsetKeyErrMsg)
	}

	return pageReq, nil
}
