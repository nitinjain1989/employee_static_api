package handlers

import (
	"encoding/json"
	"net/http"
	"static-api/data"
	"static-api/models"
)

func GetEmployees(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := models.APIResponse{
		Status:  "success",
		Message: "Employees fetched successfully",
		Data: models.EmployeeData{
			Employees: data.Employees,
		},
		Meta: models.Meta{
			TotalCount:  len(data.Employees),
			Page:        1,
			PageSize:    20,
			HasNextPage: false,
		},
	}

	json.NewEncoder(w).Encode(response)
}
