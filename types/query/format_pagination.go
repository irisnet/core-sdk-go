package query

func FormatPagination(pagination *PageRequest) *PageRequest {
	pagination.CountTotal = false
	if pagination.Limit == 0 || pagination.Limit > 100 {
		pagination.Limit = 100
	}
	return pagination
}
