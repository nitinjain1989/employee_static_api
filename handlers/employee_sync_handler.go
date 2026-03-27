package handlers

import (
	"net/http"
	"static-api/dto"
	"static-api/response"
	"static-api/services"

	"github.com/gin-gonic/gin"
)

type SyncHandler struct {
	service *services.SyncService
}

func NewSyncHandler(s *services.SyncService) *SyncHandler {
	return &SyncHandler{service: s}
}

/*
func (h *SyncHandler) SyncHandler(c *gin.Context) {

		var req dto.SyncRequest

		if err := c.ShouldBind(&req); err != nil {
			c.JSON(400, response.MessageResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		resp, err := h.service.Sync(c, req)
		if err != nil {
			c.JSON(500, response.MessageResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		c.JSON(200, resp)
	}
*/
func (h *SyncHandler) Sync(c *gin.Context) {

	var req dto.SyncRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	resp, err := h.service.Sync(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.MessageResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    resp,
	})
}
