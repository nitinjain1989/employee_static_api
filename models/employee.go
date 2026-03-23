package models

type Mobile struct {
	ID         string `json:"id,omitempty"`
	EmployeeID string `json:"employee_id,omitempty"`
	Type       string `json:"type"` // home, office, other
	Number     string `json:"number"`
}

type Employee struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Designation string `json:"designation"`
	Department  string `json:"department"`
	IsActive    bool   `json:"is_active"`
	ImgURL      string `json:"img_url"`
	Email       string `json:"email" binding:"required,email"`
	City        string `json:"city"`
	Country     string `json:"country"`
	JoiningDate string `json:"joining_date"`

	Mobiles []Mobile `json:"mobiles,omitempty"`
}

type Meta struct {
	TotalCount  int  `json:"total_count"`
	Page        int  `json:"page"`
	PageSize    int  `json:"page_size"`
	HasNextPage bool `json:"has_next_page"`
}

type EmployeeFilter struct {
	Search      string   `json:"search"`      // 🔍 name search
	Designation []string `json:"designation"` // 🎯 multi-select
	Department  []string `json:"department"`  // 🎯 multi-select
	Status      string   `json:"status"`      // active / inactive

	// 📦 Pagination (optional)
	Limit  *int `json:"limit"`
	Offset *int `json:"offset"`
}

type MessageResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type EmployeeListResponse struct {
	Status  string     `json:"status"`
	Message string     `json:"message"`
	Data    []Employee `json:"data"`
	Meta    *Meta      `json:"meta,omitempty"`
}

type EmployeeDetailResponse struct {
	Status  string    `json:"status"`
	Message string    `json:"message"`
	Data    *Employee `json:"data"`
}

type FilterOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type EmployeeFilters struct {
	Designations []string       `json:"designations"`
	Departments  []string       `json:"departments"`
	Statuses     []FilterOption `json:"statuses"`
	MobileTypes  []FilterOption `json:"mobile_types"`
}

type EmployeeFiltersResponse struct {
	Status  string           `json:"status"`
	Message string           `json:"message"`
	Data    *EmployeeFilters `json:"data"`
}
