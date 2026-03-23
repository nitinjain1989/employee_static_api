package handlers

import (
	"static-api/models"
	"static-api/services"
	"static-api/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
	service *services.EmployeeService
}

func NewEmployeeHandler(s *services.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{service: s}
}

// GetEmployees godoc
// @Summary Get employees with filters
// @Description Fetch employees with pagination, search, and filters
// @Tags employees
// @Accept json
// @Produce json
// @Param limit query int false "Number of records to fetch" example(10)
// @Param offset query int false "Number of records to skip" example(0)
// @Param search query string false "Search keyword"
// @Param status query string false "Filter by status"
// @Param designation query string false "Comma-separated designations"
// @Param department query string false "Comma-separated departments"
// @Success 200 {object} models.EmployeeListResponse
// @Failure 500 {object} models.MessageResponse "Internal server error"
// @Router /employees [get]
func (h *EmployeeHandler) GetEmployees(c *gin.Context) {

	filter := models.EmployeeFilter{
		Limit:  utils.GetInt(c.Query("limit")),
		Offset: utils.GetInt(c.Query("offset")),
		Search: c.Query("search"),
		Status: c.Query("status"),
	}

	if d := c.Query("designation"); d != "" {
		filter.Designation = strings.Split(d, ",")
	}

	if d := c.Query("department"); d != "" {
		filter.Department = strings.Split(d, ",")
	}

	employees, meta, err := h.service.GetEmployees(c.Request.Context(), filter)

	if err != nil {
		c.JSON(500, models.MessageResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(200, models.EmployeeListResponse{
		Status:  "success",
		Message: "Employees fetched successfully",
		Data:    employees,
		Meta:    &meta,
	})
}

// GetEmplyoeeByID godoc
// @Summary Get employee by ID
// @Description Fetch a single employee by its unique ID
// @Tags employees
// @Accept json
// @Produce json
// @Param id path string true "Employee ID"
// @Success 200 {object} models.EmployeeDetailResponse
// @Failure 400 {object} models.MessageResponse "Missing or invalid ID"
// @Failure 500 {object} models.MessageResponse "Internal server error"
// @Router /employees/{id} [get]
func (eh *EmployeeHandler) GetEmplyoeeByID(c *gin.Context) {

	id := c.Param("id")

	if id == "" {
		c.JSON(400, models.MessageResponse{
			Status:  "error",
			Message: "Missing id",
		})
		return
	}

	employee, err := eh.service.GetEmployeeByID(id)

	if err != nil {
		c.JSON(500, models.MessageResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(200, models.EmployeeDetailResponse{
		Status:  "success",
		Message: "Employee fetched successfully",
		Data:    employee,
	})
}

// CreateEmployee godoc
// @Summary Create a new employee
// @Description Create a new employee with the provided details
// @Tags employees
// @Accept json
// @Produce json
// @Param employee body models.CreateEmployeeRequest true "Employee payload"
// @Success 201 {object} models.EmployeeDetailResponse "Employee created successfully"
// @Failure 400 {object} models.MessageResponse "Invalid request body or validation error"
// @Failure 500 {object} models.MessageResponse "Internal server error"
// @Router /employees [post]
func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
	var req models.CreateEmployeeRequest

	// Parse request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, models.MessageResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	// Validation
	if req.Name == "" {
		c.JSON(400, models.MessageResponse{
			Status:  "error",
			Message: "Name is required",
		})
		return
	}

	var mobiles []models.Mobile
	for _, m := range req.Mobiles {
		mobiles = append(mobiles, models.Mobile{
			Number: m.Number,
			Type:   m.Type,
		})
	}

	employee := models.Employee{
		Name:        req.Name,
		Email:       req.Email,
		Designation: req.Designation,
		Department:  req.Department,
		City:        req.City,
		Country:     req.Country,
		ImgURL:      req.ImgURL,
		JoiningDate: req.JoiningDate,
		Mobiles:     mobiles, // ✅ important
	}

	// Call service
	id, err := h.service.CreateEmployee(employee)
	if err != nil {
		c.JSON(500, models.MessageResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	// Success response
	c.JSON(201, models.EmployeeDetailResponse{
		Status:  "success",
		Message: "Employee created successfully",
		Data: &models.Employee{
			ID: id,
		},
	})
}

// UpdateEmployee godoc
// @Summary Update an existing employee
// @Description Update employee details by ID
// @Tags employees
// @Accept json
// @Produce json
// @Param id path string true "Employee ID"
// @Param employee body models.Employee true "Employee payload"
// @Success 200 {object} models.MessageResponse "Employee updated successfully"
// @Failure 400 {object} models.MessageResponse "Missing ID or invalid request body"
// @Failure 500 {object} models.MessageResponse "Internal server error"
// @Router /api/employees/{id} [patch]
func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	id := c.Param("id") // 👈 replaces mux.Vars

	if id == "" {
		c.JSON(400, models.MessageResponse{
			Status:  "error",
			Message: "Missing id",
		})
		return
	}

	var emp models.Employee

	// Parse request body
	if err := c.ShouldBindJSON(&emp); err != nil {
		c.JSON(400, models.MessageResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	// Call service
	if err := h.service.UpdateEmployee(id, emp); err != nil {
		c.JSON(500, models.MessageResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	// Success response
	c.JSON(200, models.MessageResponse{
		Status:  "success",
		Message: "Employee updated successfully",
	})
}

// DeleteEmployee godoc
// @Summary Delete an employee
// @Description Delete an employee by its ID
// @Tags employees
// @Accept json
// @Produce json
// @Param id path string true "Employee ID"
// @Success 200 {object} models.MessageResponse "Employee deleted successfully"
// @Failure 400 {object} models.MessageResponse "Missing employee ID"
// @Failure 500 {object} models.MessageResponse "Internal server error"
// @Router /employees/{id} [delete]
func (h *EmployeeHandler) DeleteEmployee(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(400, models.MessageResponse{
			Status:  "error",
			Message: "Employee ID is required",
		})
		return
	}

	err := h.service.DeleteEmployee(id)
	if err != nil {
		c.JSON(500, models.MessageResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(200, models.MessageResponse{
		Status:  "success",
		Message: "Employee deleted successfully",
	})
}
