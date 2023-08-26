package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/chunnior/geo/internal/handlers"
	"github.com/chunnior/geo/internal/handlers/middlewares"
	"github.com/chunnior/geo/internal/util"
	"github.com/chunnior/geo/internal/util/log"
)

func main() {
	cfg, err := util.ReadConfig(util.GetEnvironment())
	if err != nil {
		log.Panic("Fatal error loading config: ", err.Error())
	}

	app := fiber.New(fiber.Config{
		ReadBufferSize: 7200, // Allow bigger headers
	})

	app.Use(logger.New())

	nrProvider := middlewares.InitializeNewRelicProvider()
	// Apply in all environment except local
	if !util.IsLocal(util.GetEnvironment()) {
		middlewares.ConfigureNewRelic(nrProvider, app)
	}

	// Apply only in development environment
	if util.IsDevelopmentEnvironment() {
		middlewares.ConfigureSwagger(app)
	}

	handlers.InitRoutes(app, *cfg, nrProvider)

	log.Info("Server listening on port 8080")
	_ = app.Listen(":8080")
}
