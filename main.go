package main

import (
	"embed"
	"net/http"
	"umineko_quote/internal/controllers"
	"umineko_quote/internal/quote"
	"umineko_quote/internal/routes"
	"umineko_quote/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	app := fiber.New()

	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} ${path} ${queryParams}\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))

	quoteService := quote.NewService()
	service := controllers.NewService(quoteService)
	routes.PublicRoutes(service, app)

	app.Use("/", filesystem.New(filesystem.Config{
		Root:       http.FS(staticFiles),
		PathPrefix: "static",
		Browse:     false,
	}))

	utils.StartServerWithGracefulShutdown(app, ":3000")
}
