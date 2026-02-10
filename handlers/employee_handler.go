package handlers

import (
	"encoding/json"
	"net/http"
	"static-api/config"
	"static-api/models"
)

func GetEmployees(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	req, err := config.NewSupabaseRequest("GET", "/employees", nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer resp.Body.Close()

	var employees []models.Employee
	if err := json.NewDecoder(resp.Body).Decode(&employees); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(models.APIResponse{
		Status:  "success",
		Message: "Employees fetched successfully",
		Data: models.EmployeeData{
			Employees: employees,
		},
		Meta: models.Meta{
			TotalCount: len(employees),
		},
	})
}

func CreateEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var emp models.Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	body, _ := json.Marshal(emp)

	req, err := config.NewSupabaseRequest("POST", "/employees", body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		http.Error(w, "Failed to insert employee", resp.StatusCode)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
	})
}
