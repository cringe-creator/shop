package handlers

import (
	"net/http"

	"shop/internal/services"

	"github.com/gin-gonic/gin"
)

type PurchaseHandler struct {
	PurchaseService *services.PurchaseService
}

func (h *PurchaseHandler) BuyItem(c *gin.Context) {
	var req struct {
		ItemName string `json:"item_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	err := h.PurchaseService.BuyItem(userID.(int), req.ItemName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "purchase successful"})
}
