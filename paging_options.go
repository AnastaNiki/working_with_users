package working_with_users

type Paging_options struct {
	PageNumber int    `json:"pageNumber"`
	PageSize   int    `json:"pageSize"`
	SortField  string `json:"sortField"`
	Descending bool   `json:"descending"`
}

/*type Paging_options struct {
	PageNumber int    `json:"pageNumber" binding:"required"`
	PageSize   int    `json:"pageSize" binding:"required"`
	SortField  string `json:"sortField" binding:"required"`
	Descending bool   `json:"descending" binding:"required"`
}*/
