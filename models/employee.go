package models

type Location struct {
	City    string `json:"city"`
	Country string `json:"country"`
}

type Employee struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Designation string   `json:"designation"`
	Department  string   `json:"department"`
	IsActive    bool     `json:"is_active"`
	ImgURL      string   `json:"img_url"`
	Email       string   `json:"email"`
	Location    Location `json:"location"`
	JoiningDate string   `json:"joining_date"`
}

type EmployeeData struct {
	Employees []Employee `json:"employees"`
}

type Meta struct {
	TotalCount  int  `json:"total_count"`
	Page        int  `json:"page"`
	PageSize    int  `json:"page_size"`
	HasNextPage bool `json:"has_next_page"`
}

type APIResponse struct {
	Status  string       `json:"status"`
	Message string       `json:"message"`
	Data    EmployeeData `json:"data"`
	Meta    Meta         `json:"meta"`
}
