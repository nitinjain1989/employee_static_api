package dto

type MobileResponse struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Number string `json:"number"`
}

type EmployeeResponse struct {
	ID          string           `json:"id,omitempty"` // optional: keep or remove
	Name        string           `json:"name"`
	Designation string           `json:"designation"`
	Department  string           `json:"department"`
	IsActive    bool             `json:"is_active"`
	ImgURL      string           `json:"img_url"`
	Email       string           `json:"email"`
	City        string           `json:"city"`
	Country     string           `json:"country"`
	JoiningDate string           `json:"joining_date"`
	UpdatedAt   string           `json:"updated_at"`
	DeletedAt   string           `json:"deleted_at"`
	Version     int              `json:"version"`
	CreatedAt   string           `json:"created_at"`
	Mobiles     []MobileResponse `json:"mobiles"`
}

type EmployeeListResponse struct {
	Status  string             `json:"status"`
	Message string             `json:"message"`
	Data    []EmployeeResponse `json:"data"`
	Meta    *Meta              `json:"meta,omitempty"`
}

type Meta struct {
	TotalCount       int   `json:"total_count"`
	Page             int   `json:"page"`
	PageSize         int   `json:"page_size"`
	HasNextPage      bool  `json:"has_next_page"`
	LatestUpdatedSeq int64 `json:"latest_updated_seq"`
}

type EmployeeDetailResponse struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Data    *EmployeeResponse `json:"data"`
}
