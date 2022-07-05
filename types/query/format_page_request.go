package query

import sdk "github.com/irisnet/core-sdk-go/types"

const (
	CountTotalErrMsg = "pageRequest error: CountTotal is not supported, must be false"
	LimitErrMsg      = "pageRequest error: Limit cannot be empty and cannot exceed 100"
)

func FormatPageRequest(pageReq *PageRequest) (*PageRequest, sdk.Error) {
	if pageReq.CountTotal {
		return pageReq, sdk.Wrapf(CountTotalErrMsg)
	}
	pageReq.CountTotal = false
	if pageReq.Limit == 0 || pageReq.Limit > 100 {
		return pageReq, sdk.Wrapf(LimitErrMsg)
	}

	return pageReq, nil
}
