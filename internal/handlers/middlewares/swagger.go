package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func ConfigureSwagger(app *fiber.App) {

	app.Get("/swagger/*", swagger.HandlerDefault)     // default
	app.Get("/swagger/*", swagger.New(swagger.Config{ // custom
		URL:          "http://example.com/doc.json",
		DeepLinking:  false,
		DocExpansion: "none",
	}))
}
