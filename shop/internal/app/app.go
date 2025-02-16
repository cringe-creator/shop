package app

import (
	"log"
	"shop/internal/db"
	"shop/internal/handlers"

	"github.com/gin-gonic/gin"
)

type App struct {
	Router *gin.Engine
	DB     *db.Database
}

func NewApp() *App {
	app := &App{}

	app.DB = db.NewDatabase()

	app.Router = gin.Default()
	handlers.RegisterRoutes(app.Router, app.DB.DB)

	return app
}

func (a *App) Run() {
	log.Println("Server is running on port 8080")
	a.Router.Run(":8080")
}
