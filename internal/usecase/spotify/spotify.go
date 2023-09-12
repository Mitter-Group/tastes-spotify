package spotify

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/big"

	"github.com/chunnior/spotify/internal/integration/external"
	"github.com/chunnior/spotify/internal/models"
	configRepo "github.com/chunnior/spotify/internal/repository/data"
	"github.com/chunnior/spotify/internal/util/log"
	"github.com/chunnior/spotify/pkg/aws/dynamodb"
	"github.com/valyala/fasthttp"

	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

type Implementation struct {
	repo            configRepo.Spec
	extSpotify      external.Integration
	spotifyAuthConf models.SpotifyAuthConf
}

func NewUseCase(spotifyAuthConf models.SpotifyAuthConf, repo configRepo.Spec, external external.Integration) *Implementation {
	return &Implementation{
		repo:            repo,
		extSpotify:      external,
		spotifyAuthConf: spotifyAuthConf,
	}
}

func (i *Implementation) Save(ctx context.Context, tasteReq models.DataRequest) (models.Data, error) {

	randInt, _ := rand.Int(rand.Reader, big.NewInt(100))

	taste := models.Data{
		UserId: tasteReq.UserId,
		Data: []models.DataDetails{
			{
				ID:   randInt.String(),
				Name: "Track 1",
			},
		},
	}

	return taste, nil
}

func (i *Implementation) Get(ctx context.Context, dataType string, userId string) (*models.Data, error) {

	switch dataType {
	case "tracks":
		log.InfofWithContext(ctx, "[GET] Getting top tracks for user: %s", userId)
		return i.getTopTracks(ctx, userId)
	case "artists":
		log.InfofWithContext(ctx, "[GET] Getting top artists for user: %s", userId)
		return nil, nil
	case "genres":
		log.InfofWithContext(ctx, "[GET] Getting top genres for user: %s", userId)
		return nil, nil
	default:
		return nil, nil
	}

}

func (i *Implementation) getTopTracks(ctx context.Context, userId string) (*models.Data, error) {

	// TODO: Implementar llamada a API de Spotify
	// Debemos buscar la data en memoria o base de datos, si no existe, llamar a la API de Spotify

	data, err := i.repo.GetData(userId, "tracks")
	if err != nil {
		if err == dynamodb.ErrNotFound {
			log.Errorf("[GET] Data not found for user: %s", userId)

			tracks, err := i.extSpotify.GetTopTracks(ctx)

			if err != nil {
				log.Errorf("[GET] Error getting top tracks for user: %s", userId)
				return nil, err
			}

			log.Info(tracks)

			return nil, nil
		}
		return nil, err
	}

	if data != nil {
		return data, nil
	}

	return data, nil
}

func (i *Implementation) Login(ctx context.Context) (*string, string, error) {
	state := generateRandomState()

	spotifyAuthConf := i.spotifyAuthConf.GetSpotifyAuthConf()

	var auth = spotifyauth.New(
		spotifyauth.WithRedirectURL(i.spotifyAuthConf.RedirectURI),
		spotifyauth.WithScopes(spotifyauth.ScopeUserTopRead),
		spotifyauth.WithClientID(spotifyAuthConf.ClientID),
		spotifyauth.WithClientSecret(spotifyAuthConf.ClientSecret),
	)

	authURL := auth.AuthURL(state)
	log.Info(authURL)
	return &authURL, state, nil
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

func (i *Implementation) HandleCallback(ctx context.Context, code string) error {

	spotifyAuthConf := i.spotifyAuthConf.GetSpotifyAuthConf()

	args := fasthttp.Args{}
	args.Add("grant_type", "authorization_code")
	args.Add("code", code)
	args.Add("redirect_uri", i.spotifyAuthConf.RedirectURI)
	args.Add("client_id", spotifyAuthConf.ClientID)
	args.Add("client_secret", spotifyAuthConf.ClientSecret)

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
		return err
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		return errors.New("Spotify returned an error: " + string(resp.Body()))
	}

	// 	/*
	// 			{
	// 		    "access_token": "token",
	// 		    "token_type": "Bearer",
	// 		    "expires_in": 3600,
	// 		    "refresh_token": "refreshtoken",
	// 		    "scope": "user-top-read"
	// 			}
	// 	*/

	var tokenResponse struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		Scope        string `json:"scope"`
	}

	err = json.Unmarshal(resp.Body(), &tokenResponse)
	if err != nil {
		return err
	}

	// TODO:
	// Se debe guardar el token en memoria con un ttl segun el expires_in
	// Se debe guardar el refresh token en base de datos

	// return &oauth2.Token{
	// 	AccessToken: tokenResponse.AccessToken,
	// 	TokenType:   tokenResponse.TokenType,
	// }, nil

	return nil
}
