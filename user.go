package working_with_users

// по тз (чтобы была возможность создать ошибку правильного типа)
type User struct {
	Id         int    `json:"id" db:"id"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Patronymic string `json:"patronymic"`
	Login      string `json:"login"`
	Archive    bool   `json:"archive"`
}

type UserPaging struct {
	FirstName  string        `json:"firstName"`
	LastName   string        `json:"lastName"`
	Patronymic string        `json:"patronymic"`
	Login      string        `json:"login"`
	Archive    bool          `json:"archive"`
	Paging     PagingOptions `json:"pagingOptions"`
}

type PagingOptions struct {
	PageNumber int
	PageSize   int    `json:"pageSize"`
	SortField  string `json:"sortField"`
	Descending bool   `json:"descending"`
}

type PagingResult struct {
	PageNumber  int   `json:"pageNumber"`
	PageSize    int   `json:"pageSize"`
	PageTotal   int64 `json:"pageTotal"`
	ItemsCount  int64 `json:"itemsCount"`
	HasNextPage bool  `json:"hasNextPage"`
}

//структура для автоматической проверки обязательных полей
/*type User struct {
	Id         int    `json:"-" db:"id"`
	FirstName  string `json:"firstName" binding:"required"`
	LastName   string `json:"lastName" binding:"required"`
	Patronymic string `json:"patronymic"`
	Login      string `json:"login" binding:"required"`
	Archive    bool   `json:"archive"`
}*/
