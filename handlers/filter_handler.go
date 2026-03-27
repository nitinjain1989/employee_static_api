package handlers

import (
	"static-api/dto"
	"static-api/response"

	"github.com/gin-gonic/gin"
)

// GetEmployeeFilters godoc
// @Summary Get employee filter options
// @Description Fetch available filter values like designation, department, status, etc.
// @Tags employees
// @Accept json
// @Produce json
// @Success 200 {object} dto.EmployeeFiltersResponse "Filter data fetched successfully"
// @Failure 500 {object} response.MessageResponse "Internal server error"
// @Router /employees/filters [get]
func (h *EmployeeHandler) GetEmployeeFilters(c *gin.Context) {
	data, err := h.service.GetEmployeeFilters()

	if err != nil {
		c.JSON(500, response.MessageResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(200, dto.EmployeeFiltersResponse{
		Status:  "success",
		Message: "Filters fetched successfully",
		Data:    data,
	})
}
