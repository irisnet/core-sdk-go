package query

func FormatPageRequest(pageReq *PageRequest) *PageRequest {
	pageReq.CountTotal = false
	if pageReq.Limit == 0 || pageReq.Limit > 100 {
		pageReq.Limit = 100
	}
	return pageReq
}
