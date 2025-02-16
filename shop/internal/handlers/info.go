package handlers

import (
	"net/http"
	"shop/internal/services"

	"github.com/gin-gonic/gin"
)

type InfoHandler struct {
	InfoService *services.InfoService
}

func (h *InfoHandler) GetUserInfo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userInfo, err := h.InfoService.GetUserInfo(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user data"})
		return
	}

	c.JSON(http.StatusOK, userInfo)
}
