package handlers

import (
	"encoding/json"
	"net/http"
	"static-api/config"
	"static-api/models"
)

func GetEmployees(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	client := config.NewSupabaseClient()

	var employees []models.Employee
	err := client.DB.
		From("employees").
		Select("*").
		Execute(&employees)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	response := models.APIResponse{
		Status:  "success",
		Message: "Employees fetched successfully",
		Data: models.EmployeeData{
			Employees: employees,
		},
		Meta: models.Meta{
			TotalCount: len(employees),
		},
	}

	json.NewEncoder(w).Encode(response)
}

func CreateEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var emp models.Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	client := config.NewSupabaseClient()

	_, err := client.DB.
		From("employees").
		Insert(emp).
		Execute(nil)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
	})
}
