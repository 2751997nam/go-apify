package types

type Meta struct {
	HasNext    bool  `json:"has_next"`
	PageId     int   `json:"page_id"`
	PageSize   int   `json:"page_size"`
	PageCount  int   `json:"page_count"`
	TotalCount int64 `json:"total_count"`
}
