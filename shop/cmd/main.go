package main

import (
	"shop/internal/app"
)

func main() {
	application := app.NewApp()
	application.Run()
}
