package handlers

type PaginatedResponse[T any] struct {
	Data       []T `json:"data"`
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalCount int `json:"totalCount"`
	TotalPages int `json:"totalPages"`
}

// calcTotalPages calculates the total pages of paginated data available (equivalent to the last page number)
func calcTotalPages(totalCount, pageSize int) int {
	return (totalCount + pageSize - 1) / pageSize // round up trick for totalCount/pageSize
}
