package entity

type PaginateSearch struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Search   string `json:"search"`
}

type PaginatedResponse struct {
	CurrentPage int64       `json:"current_page,omitempty"`
	TotalPage   int         `json:"total_page,omitempty"`
	TotalItems  int64       `json:"total_items,omitempty"`
	Data        interface{} `json:"data,omitempty"`
}

func ToPaginatedResponse(currentPage int64, totalPage int, totalItems int64, data interface{}) *PaginatedResponse {
	return &PaginatedResponse{
		CurrentPage: currentPage,
		TotalPage:   totalPage,
		TotalItems:  totalItems,
		Data:        data,
	}
}
