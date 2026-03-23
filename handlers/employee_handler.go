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
		c.JSON(500, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, models.APIResponse{
		Status:  "success",
		Message: "Employees fetched successfully",
		Data: models.EmployeeData{
			Employees: employees,
		},
		Meta: &meta,
	})
}

/*func (eh *EmployeeHandler) GetEmployees(c *gin.Context) {

	employees, meta, err := eh.service.GetEmployees(c)

	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, models.APIResponse{
		Status:  "success",
		Message: "Employees fetched successfully",
		Data:    models.EmployeeData{Employees: employees},
		Meta:    &meta,
	})
}*/

func (eh *EmployeeHandler) GetEmplyoeeByID(c *gin.Context) {

	id := c.Param("id")

	if id == "" {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Missing id",
		})
		return
	}

	employee, err := eh.service.GetEmployeeByID(id)

	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, models.APIResponse{
		Status:  "success",
		Message: "Employee fetched successfully",
		Data: models.EmployeeData{
			Employee: employee,
		},
	})
}

func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
	var emp models.Employee

	// Parse request body
	if err := c.ShouldBindJSON(&emp); err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Validation
	if emp.Name == "" {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Name is required",
		})
		return
	}

	// Call service
	id, err := h.service.CreateEmployee(emp)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Success response
	c.JSON(201, models.APIResponse{
		Status:  "success",
		Message: "Employee created successfully",
		Data: models.EmployeeData{
			Employee: &models.Employee{
				ID: id,
			},
		},
	})
}

func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	id := c.Param("id") // 👈 replaces mux.Vars

	if id == "" {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Missing id",
		})
		return
	}

	var emp models.Employee

	// Parse request body
	if err := c.ShouldBindJSON(&emp); err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Call service
	if err := h.service.UpdateEmployee(id, emp); err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Success response
	c.JSON(200, models.APIResponse{
		Status:  "success",
		Message: "Employee updated successfully",
	})
}

func (h *EmployeeHandler) DeleteEmployee(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Employee ID is required",
		})
		return
	}

	err := h.service.DeleteEmployee(id)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, models.APIResponse{
		Status:  "success",
		Message: "Employee deleted successfully",
	})
}
