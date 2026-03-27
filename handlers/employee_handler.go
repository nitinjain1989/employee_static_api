package handlers

import (
	"errors"
	"static-api/dto"
	"static-api/response"
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
// @Success 200 {object} dto.EmployeeListResponse
// @Failure 500 {object} response.MessageResponse "Internal server error"
// @Router /employees [get]
func (h *EmployeeHandler) GetEmployees(c *gin.Context) {

	filter := dto.EmployeeFilterRequest{
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
		c.JSON(500, response.MessageResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(200, dto.EmployeeListResponse{
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
// @Success 200 {object} dto.EmployeeDetailResponse
// @Failure 400 {object} response.MessageResponse "Missing or invalid ID"
// @Failure 500 {object} response.MessageResponse "Internal server error"
// @Router /employees/{id} [get]
func (eh *EmployeeHandler) GetEmplyoeeByID(c *gin.Context) {

	id := c.Param("id")

	if id == "" {
		c.JSON(400, response.MessageResponse{
			Status:  "error",
			Message: "Missing id",
		})
		return
	}

	employee, err := eh.service.GetEmployeeByID(id)

	if err != nil {
		c.JSON(500, response.MessageResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(200, dto.EmployeeDetailResponse{
		Status:  "success",
		Message: "Employee fetched successfully",
		Data:    employee,
	})
}

// CreateEmployee godoc
// @Summary Create a new employee
// @Description Creates a new employee with basic details and up to 3 mobile numbers. Validates required fields before creation.
// @Tags employees
// @Accept json
// @Produce json
// @Param employee body dto.CreateEmployeeRequest true "Employee payload (max 3 mobiles allowed)"
// @Success 201 {object} dto.EmployeeDetailResponse "Employee created successfully"
// @Failure 400 {object} response.MessageResponse "Invalid request body, missing name, or too many mobiles"
// @Failure 500 {object} response.MessageResponse "Internal server error"
// @Router /employees [post]
func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
	var req dto.CreateEmployeeRequest

	// Parse request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, response.MessageResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	// Validation
	if req.Name == "" {
		c.JSON(400, response.MessageResponse{
			Status:  "error",
			Message: "Name is required",
		})
		return
	}

	if len(req.Mobiles) > 3 {
		c.JSON(400, response.MessageResponse{
			Status:  "error",
			Message: "maximum 3 mobiles allowed per employee",
		})
		return
	}

	// Call service
	employee, err := h.service.CreateEmployee(req)
	if err != nil {
		c.JSON(500, response.MessageResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	// Success response
	c.JSON(201, dto.EmployeeDetailResponse{
		Status:  "success",
		Message: "Employee created successfully",
		Data:    employee,
	})
}

// UpdateEmployee godoc
// @Summary Update an existing employee
// @Description Updates employee details by ID. Supports partial updates. Maximum 3 mobile numbers allowed. Returns conflict if version mismatch (last-write handling).
// @Tags employees
// @Accept json
// @Produce json
// @Param id path string true "Employee ID"
// @Param employee body dto.UpdateEmployeeRequest true "Employee update payload (max 3 mobiles allowed)"
// @Success 200 {object} dto.EmployeeDetailResponse "Employee updated successfully"
// @Failure 400 {object} response.MessageResponse "Missing ID, invalid request body, or validation error"
// @Failure 404 {object} response.MessageResponse "Employee not found"
// @Failure 409 {object} response.MessageResponse "Conflict error (e.g., version mismatch)"
// @Failure 500 {object} response.MessageResponse "Internal server error"
// @Router /api/employees/{id} [patch]
func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	id := c.Param("id") // 👈 replaces mux.Vars

	if id == "" {
		c.JSON(400, response.MessageResponse{
			Status:  "error",
			Message: "Missing id",
		})
		return
	}

	var req dto.UpdateEmployeeRequest

	// Parse request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, response.MessageResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	if len(req.Mobiles) > 3 {
		c.JSON(400, response.MessageResponse{
			Status:  "error",
			Message: "maximum 3 mobiles allowed per employee",
		})
		return
	}

	emp, err := h.service.UpdateEmployee(id, req)
	// Call service
	if err != nil {

		if errors.Is(err, utils.ErrConflict) {
			c.JSON(409, response.MessageResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		if errors.Is(err, utils.ErrNotFound) {
			c.JSON(404, response.MessageResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		c.JSON(500, response.MessageResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	// Success response
	c.JSON(200, dto.EmployeeDetailResponse{
		Status:  "success",
		Message: "Employee updated successfully",
		Data:    emp,
	})
}

// DeleteEmployee godoc
// @Summary Delete an employee
// @Description Delete an employee by its ID
// @Tags employees
// @Accept json
// @Produce json
// @Param id path string true "Employee ID"
// @Success 200 {object} response.MessageResponse "Employee deleted successfully"
// @Failure 400 {object} response.MessageResponse "Missing employee ID"
// @Failure 500 {object} response.MessageResponse "Internal server error"
// @Router /employees/{id} [delete]
func (h *EmployeeHandler) DeleteEmployee(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(400, response.MessageResponse{
			Status:  "error",
			Message: "Employee ID is required",
		})
		return
	}

	err := h.service.DeleteEmployee(id)
	if err != nil {
		c.JSON(500, response.MessageResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(200, response.MessageResponse{
		Status:  "success",
		Message: "Employee deleted successfully",
	})
}
