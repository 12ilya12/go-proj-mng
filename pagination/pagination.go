package pagination

//Входные параметры пагинации
type PagingOptions struct {
	Page     int    `json:"page,omitempty" schema:"page"`
	PageSize int    `json:"page_size,omitempty" schema:"page_size"`
	Order    Order  `json:"order,omitempty" schema:"order"`
	OrderBy  string `json:"order_by,omitempty" schema:"order_by"`
}

//Выходные параметры пагинации
type Pagination struct {
	TotalItems int64         `json:"total_items"`
	TotalPages int64         `json:"total_pages"`
	Options    PagingOptions `json:"options"`
}

//Объект-дженерик для формирования ответа на запрос с пагинацией
type Paging[T any] struct {
	Items      []T        `json:"items"`
	Pagination Pagination `json:"pagination"`
}

type Order string

const (
	Asc  Order = "asc"
	Desc Order = "desc"
)
