package main

import "github.com/yrss1/todo/internal/app"

// @title Todo API
// @version 1.0
// @description This is a sample server for a todo application.
// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	app.Run()
}
