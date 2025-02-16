package handlers

import (
	"net/http"

	"shop/internal/services"

	"github.com/gin-gonic/gin"
)

type HistoryHandler struct {
	HistoryService *services.HistoryService
}

func (h *HistoryHandler) GetUserHistory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	history, err := h.HistoryService.GetUserHistory(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch history"})
		return
	}

	c.JSON(http.StatusOK, history)
}
