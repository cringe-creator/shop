package handlers

import (
	"database/sql"
	"shop/internal/middleware"
	"shop/internal/repositories"
	"shop/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, db *sql.DB) {
	userRepo := &repositories.UserRepository{DB: db}
	transactionRepo := &repositories.TransactionRepository{DB: db}
	authService := &services.AuthService{UserRepo: userRepo}
	purchaseRepo := &repositories.PurchaseRepository{DB: db}

	transactionService := &services.TransactionService{
		TransactionRepo: transactionRepo,
		UserRepo:        userRepo,
	}
	purchaseService := &services.PurchaseService{
		UserRepo:     userRepo,
		PurchaseRepo: purchaseRepo,
	}
	historyService := &services.HistoryService{
		TransactionRepo: transactionRepo,
		PurchaseRepo:    purchaseRepo,
	}

	infoService := &services.InfoService{
		UserRepo:        userRepo,
		TransactionRepo: transactionRepo,
	}

	authHandler := &AuthHandler{AuthService: authService}
	transactionHandler := &TransactionHandler{TransactionService: transactionService}
	purchaseHandler := &PurchaseHandler{PurchaseService: purchaseService}
	historyHandler := &HistoryHandler{HistoryService: historyService}
	infoHandler := &InfoHandler{InfoService: infoService}

	r.POST("/api/auth/register", authHandler.Register)
	r.POST("/api/auth/login", authHandler.Login)

	protected := r.Group("/api")
	protected.Use(middleware.JWTMiddleware())
	protected.POST("/sendCoin", transactionHandler.SendCoins)
	protected.POST("/buy", purchaseHandler.BuyItem)
	protected.GET("/history", historyHandler.GetUserHistory)
	protected.GET("/info", infoHandler.GetUserInfo)
}
