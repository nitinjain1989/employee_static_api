package dto

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
