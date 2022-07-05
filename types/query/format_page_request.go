package query

import sdk "github.com/irisnet/core-sdk-go/types"

const (
	CountTotalErrMsg      = "pageRequest error: CountTotal must be false"
	LimitErrMsg           = "pageRequest error: Limit cannot be empty and cannot exceed 100"
	CountTotalLimitErrMsg = "pageRequest error: CountTotal must be false and Limit cannot be empty and cannot exceed 100"
)

func FormatPageRequest(pageReq *PageRequest) (*PageRequest, sdk.Error) {
	var msg string
	if pageReq.CountTotal {
		msg = CountTotalErrMsg
	}
	pageReq.CountTotal = false
	if pageReq.Limit == 0 || pageReq.Limit > 100 {
		pageReq.Limit = 100
		if msg != "" {
			msg = CountTotalLimitErrMsg
		} else {
			msg = LimitErrMsg
		}
	}

	var err sdk.Error
	if msg != "" {
		err = sdk.Wrapf(msg)
	}

	return pageReq, err
}
