package pagination

//Входные параметры пагинации
type PagingOptions struct {
	page     int    `json:"page" schema:"page"`
	pageSize int    `json:"page_size" schema:"page_size"`
	order    Order  `json:"order" schema:"order"`
	orderBy  string `json:"order_by" schema:"order_by"`
}

//Выходные параметры пагинации
type Pagination struct {
	totalItems int           `json:"total_items"`
	totalPages int           `json:"total_pages"`
	options    PagingOptions `json:"options"`
}

//Объект-дженерик для формирования ответа на запрос с пагинацией
type Paging[T any] struct {
	items      []T        `json:"items"`
	pagination Pagination `json:"pagination"`
}

type Order string

const (
	Asc  Order = "asc"
	Desc Order = "desc"
)
