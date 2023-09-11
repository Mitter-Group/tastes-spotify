package handlers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/valyala/fasthttp"
	"golang.org/x/oauth2"

	"github.com/chunnior/spotify/internal/integration/external"
	"github.com/chunnior/spotify/internal/models"
	dataRepo "github.com/chunnior/spotify/internal/repository/data"
	spotifyUC "github.com/chunnior/spotify/internal/usecase/spotify"
	configAws "github.com/chunnior/spotify/pkg/aws"
	"github.com/chunnior/spotify/pkg/aws/dynamodb"

	"github.com/zmb3/spotify/v2"
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

	spotifyUseCase := spotifyUC.NewUseCase(repository, extSpotify)

	handler := NewHandler(spotifyUseCase)

	app.Get("/:dataType/:userId", handler.Get)

	app.Get("/login", handler.Login)
	app.Get("/callback", callbackHandler)

}

func callbackHandler(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	callbackState := c.Query("state")
	storedState := c.Cookies("oauth_state")
	if callbackState != storedState {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid state parameter")
	}

	args := fasthttp.Args{}
	args.Add("grant_type", "authorization_code")
	args.Add("code", code)
	args.Add("redirect_uri", redirectURI)
	args.Add("client_id", clientID)
	args.Add("client_secret", clientSecret)

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI("https://accounts.spotify.com/api/token")
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBody(args.QueryString())

	err := fasthttp.Do(req, resp)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to exchange token")
	}

	// Verifica el código de estado de la respuesta
	if resp.StatusCode() != fasthttp.StatusOK {
		return c.Status(resp.StatusCode()).SendString("Spotify returned an error: " + string(resp.Body()))
	}

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
	}

	err = json.Unmarshal(resp.Body(), &tokenResponse)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to parse token response")
	}

	token := &oauth2.Token{
		AccessToken: tokenResponse.AccessToken,
		TokenType:   tokenResponse.TokenType,
	}

	client := spotify.New(auth.Client(c.Context(), token))

	options := []spotify.RequestOption{
		spotify.Timerange(spotify.LongTermRange),
	}

	topArtist, err := client.CurrentUsersTopArtists(c.Context(), options...)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "No se pudieron obtener los top tracks")
	}

	var artistNames string
	for _, artist := range topArtist.Artists {
		artistNames += artist.Name + "\n"
	}

	return c.SendString(artistNames)
}
