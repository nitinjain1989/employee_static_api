package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"static-api/config"
	"static-api/models"
	"strconv"
	"strings"
)

func GetEmployees(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	queries := r.URL.Query()
	limit := queries.Get("limit")
	path := buildPath(queries)

	req, err := config.NewSupabaseRequest("GET", path, nil)
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

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		http.Error(w, "supabase error", resp.StatusCode)
		return
	}

	var employees []models.Employee
	if err := json.NewDecoder(resp.Body).Decode(&employees); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	meta := parseContentRange(resp.Header.Get("content-range"), limit)

	json.NewEncoder(w).Encode(models.APIResponse{
		Status:  "success",
		Message: "Employees fetched successfully",
		Data: models.EmployeeData{
			Employees: employees,
		},
		Meta: meta,
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

func buildInFilter(key, value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}

	items := strings.Split(value, ",")
	var cleaned []string

	for _, item := range items {
		v := strings.TrimSpace(item)
		if v != "" {
			cleaned = append(cleaned, url.QueryEscape(v))
		}
	}

	if len(cleaned) == 0 {
		return ""
	}

	return "&" + key + "=in.(" + strings.Join(cleaned, ",") + ")"
}

func parseContentRange(contentRange, limitStr string) models.Meta {
	meta := models.Meta{
		TotalCount: 0,
		PageSize:   0,
		Page:       1,
	}

	if contentRange == "" {
		return meta
	}

	// Example: "0-19/100"
	parts := strings.Split(contentRange, "/")
	if len(parts) != 2 {
		return meta
	}

	rangePart := parts[0]
	totalStr := parts[1]

	total, err := strconv.Atoi(totalStr)
	if err != nil {
		return meta
	}

	rangeParts := strings.Split(rangePart, "-")
	if len(rangeParts) != 2 {
		return meta
	}

	start, err1 := strconv.Atoi(rangeParts[0])
	end, err2 := strconv.Atoi(rangeParts[1])

	if err1 != nil || err2 != nil {
		return meta
	}

	pageSize := end - start + 1
	page := 1

	if limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			page = (start / limit) + 1
		}
	}

	meta.TotalCount = total
	meta.PageSize = pageSize
	meta.Page = page

	return meta
}
func buildPath(queries url.Values) string {
	//queries := r.URL.Query()
	limit := queries.Get("limit")
	offset := queries.Get("offset")
	search := queries.Get("search")
	designation := queries.Get("designation")
	department := queries.Get("department")
	status := queries.Get("status")

	path := "/employees?select=*"

	if search != "" {
		encodedSearch := url.QueryEscape(search)
		path += "&order=name.asc,created_at.desc,id.desc"
		path += "&name=ilike.*" + encodedSearch + "*"
	} else {
		path += "&order=created_at.desc,id.desc"
	}

	// 🎯 Filters
	path += buildInFilter("designation", designation)
	path += buildInFilter("department", department)

	if status == "active" {
		path += "&is_active=eq.true"
	} else if status == "inactive" {
		path += "&is_active=eq.false"
	}

	if limit != "" {
		path += "&limit=" + limit
	}

	if offset != "" {
		path += "&offset=" + offset
	}

	return path
}
