package dto

type EmployeeFilterRequest struct {
	Search      string   `json:"search"`      // 🔍 name search
	Designation []string `json:"designation"` // 🎯 multi-select
	Department  []string `json:"department"`  // 🎯 multi-select
	Status      string   `json:"status"`      // active / inactive

	// 📦 Pagination (optional)
	Limit  *int `json:"limit"`
	Offset *int `json:"offset"`
}
