package query

import sdk "github.com/irisnet/core-sdk-go/types"

const (
	PageReqNilErrMsg = "pageRequest error: PageRequest cannot be nil"
	CountTotalErrMsg = "pageRequest error: CountTotal is not supported, must be false"
	LimitErrMsg      = "pageRequest error: Limit cannot be empty and cannot exceed 100"
)

func FormatPageRequest(pageReq *PageRequest) (*PageRequest, sdk.Error) {
	if pageReq == nil {
		return pageReq, sdk.Wrapf(PageReqNilErrMsg)
	}
	if pageReq.CountTotal {
		return pageReq, sdk.Wrapf(CountTotalErrMsg)
	}

	if pageReq.Limit == 0 || pageReq.Limit > 100 {
		return pageReq, sdk.Wrapf(LimitErrMsg)
	}

	return pageReq, nil
}
