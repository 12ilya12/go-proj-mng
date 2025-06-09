package common

type TaskFilters struct {
	StatusId   uint `json:"status_id"`
	CategoryId uint `json:"category_id"`
	UserInfo
}

type UserInfo struct {
	UserId   int
	UserRole string
}
