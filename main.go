package main

import (
	"Wallet-System-Backend/infra"
	"Wallet-System-Backend/infra/logger"
	"Wallet-System-Backend/internal/handlers"
	"Wallet-System-Backend/internal/middleware"
	"Wallet-System-Backend/utils"
	"os"

	_ "Wallet-System-Backend/docs"

	"github.com/gofiber/swagger"

	"github.com/gofiber/fiber/v2"
)

// @title Wallet System API
// @version 1.0
// @description API documentation for Wallet System backend.
// @host localhost:3000
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	logger.InitializeLogger("app.log")

	app := fiber.New()
	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Use(middleware.CORS())
	app.Use(middleware.ErrorHandlerMiddleware)

	db := infra.InitDB()
	cache := infra.NewRedisClient()

	handlers.NewLoginHandler(app, db)

	app.Use(middleware.JWTMiddleware())

	handlers.NewProductHandler(app, db, cache)
	handlers.NewWalletHandler(app, db, cache)
	handlers.NewTransactionHandler(app, db, cache)

	utils.LoadEnv()
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	app.Listen(":" + port)
}
