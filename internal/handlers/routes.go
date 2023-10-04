package handlers

import (
	"github.com/chunnior/spotify/internal/integration/external"
	"github.com/chunnior/spotify/internal/models"
	dataRepo "github.com/chunnior/spotify/internal/repository/data"
	spotifyUC "github.com/chunnior/spotify/internal/usecase/spotify"
	configAws "github.com/chunnior/spotify/pkg/aws"
	"github.com/chunnior/spotify/pkg/aws/dynamodb"
	"github.com/gofiber/fiber/v2"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// @title          Go template
// @version        1.0
// @description    Go template API - Example for swagger
// @contact.name   Your team
// @contact.email  your-team@taste.la
// @host           go-template.dev.taste.la
// @BasePath       /
func InitRoutes(app *fiber.App, cfg models.Config, nrProvider *newrelic.Application) {

	awsConfig, _ := configAws.GetConfig(&cfg.AWS.Credentials, cfg.Env == "local")

	dynamo := dynamodb.NewDynamoClient(awsConfig,
		dynamodb.WithTable(cfg.AWS.UserSpotifyData),
	)

	repository := dataRepo.NewConfig(dynamo, cfg)

	extSpotify := external.NewIntegration(cfg)

	spotifyUseCase := spotifyUC.NewUseCase(cfg.App, cfg.Spotify, repository, extSpotify)

	handler := NewHandler(spotifyUseCase)

	app.Get("/:dataType/:userId", handler.Get)

	app.Get("/login", handler.Login)
	app.Get("/callback", handler.Callback)

}
