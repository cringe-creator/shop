package handlers

import (
	"net/http"

	"shop/internal/services"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	TransactionService *services.TransactionService
}

func (h *TransactionHandler) SendCoins(c *gin.Context) {
	var req struct {
		ReceiverName string `json:"toUser"`
		Amount       int    `json:"amount"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	senderID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	receiver, err := h.TransactionService.UserRepo.GetUserByUsername(req.ReceiverName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err = h.TransactionService.SendCoins(senderID.(int), receiver.ID, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "coins sent successfully"})
}
