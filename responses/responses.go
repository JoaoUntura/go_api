package responses

type APIResponse struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
}

type APIResponsePagination struct {
	Success    bool  `json:"success"`
	Data       any   `json:"data"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"totalItems"`
	TotalPages int   `json:"totalPages"`
}
