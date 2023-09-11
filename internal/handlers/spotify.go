package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/chunnior/spotify/internal/models"
	spotifyUC "github.com/chunnior/spotify/internal/usecase/spotify"
	"github.com/chunnior/spotify/internal/util/log"
	"github.com/chunnior/spotify/pkg/tracing"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

const (
	clientID     = "f67e9675502e409d8a85fe2e6ff07739"
	clientSecret = "609fa7897dbf41cab60624d1a840c40a"
	redirectURI  = "http://localhost:8080/callback"
)

var auth = spotifyauth.New(
	spotifyauth.WithRedirectURL(redirectURI),
	spotifyauth.WithScopes(spotifyauth.ScopeUserTopRead),
	spotifyauth.WithClientID(clientID),
	spotifyauth.WithClientSecret(clientSecret),
)

type Handler struct {
	useCase spotifyUC.UseCase
}

func NewHandler(useCase spotifyUC.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) Save(c *fiber.Ctx) error {
	ctx := startTracing(c, "saveTaste")
	defer endTracing(ctx)

	log.InfofWithContext(ctx, "[TASTE_POST] Creating Taste")

	var req models.DataRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return respondWithError(c, http.StatusBadRequest, "invalid_body")
	}

	taste, err := h.useCase.Save(ctx, req)
	if err != nil {
		return respondWithError(c, http.StatusInternalServerError, "internal_error")
	}

	log.InfofWithContext(ctx, "[TASTE_POST] Taste created: %s", taste.Data)
	return c.Status(http.StatusOK).JSON(fiber.Map{"data": taste})
}

func (h *Handler) Get(c *fiber.Ctx) error {
	ctx := startTracing(c, "getData")
	defer endTracing(ctx)

	dataType, userId := c.Params("dataType"), c.Params("userId")
	if dataType == "" || userId == "" {
		return respondWithError(c, http.StatusBadRequest, "invalid_params")
	}

	log.InfofWithContext(ctx, "[GETDATA] Getting data %s for user %s", dataType, userId)

	res, err := h.useCase.Get(ctx, dataType, userId)
	if err != nil {
		return respondWithError(c, http.StatusInternalServerError, err.Error())
	}
	if res == nil {
		return respondWithError(c, http.StatusNotFound, "data not found")
	}

	return c.Status(http.StatusOK).JSON(res)
}

func (h *Handler) Login(c *fiber.Ctx) error {
	state := generateRandomState()
	c.Cookie(&fiber.Cookie{
		Name:  "oauth_state",
		Value: state,
	})
	return c.Redirect(auth.AuthURL(state))
}

func (h *Handler) Callback(c *fiber.Ctx) error {
	if err := validateOAuthCallback(c); err != nil {
		return err
	}

	token, err := exchangeCodeForToken(c)
	if err != nil {
		return err
	}

	/*
			{
		    "access_token": "token",
		    "token_type": "Bearer",
		    "expires_in": 3600,
		    "refresh_token": "refreshtoken",
		    "scope": "user-top-read"
			}
	*/

	// TODO:
	// Se debe guardar el token en memoria con un ttl segun el expires_in
	// Se debe guardar el refresh token en base de datos

	/**** ESTO ES UNA PRUEBA DEL TOKEN ******/
	client := spotify.New(auth.Client(c.Context(), token))
	topArtist, err := fetchTopArtists(c, client)
	if err != nil {
		return err
	}

	var artistNames string
	for _, artist := range topArtist.Artists {
		artistNames += artist.Name + "\n"
	}

	return c.SendString(artistNames)
	/*****   FIN DE LA PRUEBA, SE DEBE DEVOLVER EL TOKEN *****/
}

func startTracing(c *fiber.Ctx, segmentName string) context.Context {
	ctx := tracing.CreateContextWithTransaction(c)
	txn := tracing.GetTransactionFromContext(ctx)
	txn.StartSegment(segmentName)
	return ctx
}

func endTracing(ctx context.Context) {
	txn := tracing.GetTransactionFromContext(ctx)
	txn.End()
}

func respondWithError(c *fiber.Ctx, status int, msg string) error {
	log.ErrorfWithContext(c.Context(), msg)
	return c.Status(status).JSON(formatResponse(status, msg))
}

func validateOAuthCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	callbackState, storedState := c.Query("state"), c.Cookies("oauth_state")
	if callbackState != storedState {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid state parameter")
	}
	return nil
}

func exchangeCodeForToken(c *fiber.Ctx) (*oauth2.Token, error) {
	args := fasthttp.Args{}
	args.Add("grant_type", "authorization_code")
	args.Add("code", c.Query("code"))
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
		return nil, c.Status(fiber.StatusInternalServerError).SendString("Failed to exchange token")
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		return nil, c.Status(resp.StatusCode()).SendString("Spotify returned an error: " + string(resp.Body()))
	}

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
	}

	err = json.Unmarshal(resp.Body(), &tokenResponse)
	if err != nil {
		return nil, c.Status(fiber.StatusInternalServerError).SendString("Failed to parse token response")
	}

	return &oauth2.Token{
		AccessToken: tokenResponse.AccessToken,
		TokenType:   tokenResponse.TokenType,
	}, nil
}

func fetchTopArtists(c *fiber.Ctx, client *spotify.Client) (*spotify.FullArtistPage, error) {
	options := []spotify.RequestOption{
		spotify.Timerange(spotify.LongTermRange),
	}

	topArtist, err := client.CurrentUsersTopArtists(c.Context(), options...)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "No se pudieron obtener los top tracks")
	}

	return topArtist, nil
}

func formatResponse(status int, msg string) map[string]interface{} {
	return map[string]interface{}{
		"status":  status,
		"message": msg,
	}
}

func generateRandomState() string {
	byteSize := 32
	buf := make([]byte, byteSize)
	_, err := rand.Read(buf)
	if err != nil {
		return "state-test"
	}
	return hex.EncodeToString(buf)
}
