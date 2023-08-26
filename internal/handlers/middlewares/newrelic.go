package middlewares

import (
	"os"

	"github.com/erkanzileli/nrfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/newrelic/go-agent/v3/newrelic"

	"github.com/chunnior/geo/internal/util/log"
)

func InitializeNewRelicProvider() *newrelic.Application {
	nrapp, err := newrelic.NewApplication(
		newrelic.ConfigFromEnvironment(),
		newrelic.ConfigDebugLogger(os.Stdout),
		newrelic.ConfigAppLogForwardingEnabled(true),
		newrelic.ConfigAppLogEnabled(true),
	)
	if err != nil {
		log.Error("error initializing newrelic client", err.Error())
		return nil
	}

	return nrapp
}

func ConfigureNewRelic(provider *newrelic.Application, app *fiber.App) {
	if provider != nil {
		app.Use(nrfiber.Middleware(provider))
	}
}

type ContextTag string

func InjectNewRelicTracing(name string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Locals("validator", "test")
		return c.Next()
	}
}
