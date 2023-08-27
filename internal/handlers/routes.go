package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/newrelic/go-agent/v3/newrelic"

	"github.com/chunnior/spotify/internal/entity"
	"github.com/chunnior/spotify/internal/usecase/tastes"
	"github.com/chunnior/spotify/pkg/oauth"
)

// @title          Go template
// @version        1.0
// @description    Go template API - Example for swagger
// @contact.name   Your team
// @contact.email  your-team@taste.la
// @host           go-template.dev.taste.la
// @BasePath       /
func InitRoutes(app *fiber.App, cfg entity.Config, nrProvider *newrelic.Application) {

	tasteUseCase := tastes.NewUseCase()

	handler := NewHandler(tasteUseCase)

	//oauth.Initialize() // TODO: Configurar OAUTH
	// Without token
	app.Get("/liveness", ping)
	app.Get("/readiness", ping)
	app.Post("/tokeninfo", oauth.Protected, tokeninfo)

	// api/* with auth token
	api := app.Group("/api", oauth.Protected)
	v1 := api.Group("/v1")

	// V1 Routes
	v1.Post("/auth_required_route", ping)

	app.Post("/save", handler.Save)
}
