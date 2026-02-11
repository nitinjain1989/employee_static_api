package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"static-api/config"
	"static-api/models"
)

func GetEmployees(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	req, err := config.NewSupabaseRequest("GET", "/employees?order=created_at.desc", nil)
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
		body, _ := io.ReadAll(resp.Body)
		http.Error(w, string(body), resp.StatusCode)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
	})
}

func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing employee ID", 400)
		return
	}

	var emp models.Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	body, _ := json.Marshal(emp)

	req, err := config.NewSupabaseRequest(
		"PATCH",
		"/employees?id=eq."+id,
		body,
	)
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
		body, _ := io.ReadAll(resp.Body)
		http.Error(w, string(body), resp.StatusCode)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status": "updated",
	})
}
