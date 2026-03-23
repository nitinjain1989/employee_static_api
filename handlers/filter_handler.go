package handlers

import "github.com/gin-gonic/gin"

func (h *EmployeeHandler) GetEmployeeFilters(c *gin.Context) {
	data, err := h.service.GetEmployeeFilters()

	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":   data,
	})
}
