package spotify

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"strconv"
	"time"

	"github.com/chunnior/spotify/internal/integration/external"
	"github.com/chunnior/spotify/internal/models"
	configRepo "github.com/chunnior/spotify/internal/repository/data"
	"github.com/chunnior/spotify/internal/util/log"
	"github.com/chunnior/spotify/pkg/aws/dynamodb"
	"github.com/chunnior/spotify/pkg/cache"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"

	"github.com/zmb3/spotify/v2"
)

type Implementation struct {
	repo            configRepo.Spec
	extSpotify      external.Integration
	spotifyAuthConf models.SpotifyAuthConf
	appConfig       models.AppConfig
	cacheTokens     cache.Spec
	cacheStates     cache.Spec
}

type CustomUser struct {
	*spotify.PrivateUser
}

const (
	Spotify     = "SPOTIFY"
	MemoryCache = "tokens"
	CacheSize   = 2000 // 2000 items
	CacheExpiry = 3540 // 59 minutos
)

func NewUseCase(appConfig models.AppConfig, spotifyAuthConf models.SpotifyAuthConf, repo configRepo.Spec, external external.Integration) *Implementation {
	return &Implementation{
		repo:            repo,
		extSpotify:      external,
		spotifyAuthConf: spotifyAuthConf,
		appConfig:       appConfig,
		cacheTokens:     cache.NewMemoryCache(MemoryCache, CacheSize, CacheExpiry, true),
		cacheStates:     cache.NewMemoryCache(MemoryCache, CacheSize, CacheExpiry, true),
	}
}

func (i *Implementation) Save(ctx context.Context, tasteReq models.DataRequest) (models.Data, error) {
	randInt, _ := rand.Int(rand.Reader, big.NewInt(100))

	taste := models.Data{
		UserId: tasteReq.UserId,
		Data: []models.DataDetails{
			models.TrackDetails{
				ID:          randInt.String(),
				Name:        "Track 1",
				Artists:     []models.Artist{{Name: "Artist 1", ID: "1"}},
				ReleaseDate: "2021",
			},
		},
	}

	return taste, nil
}

// Get retrieves the top data (Tracks, Artists, Genres) for a specific user.
func (i *Implementation) Get(ctx context.Context, dataType string, userId string) (*models.Data, error) {
	log.InfofWithContext(ctx, "[GET] Getting top %s for user: %s", dataType, userId)

	switch dataType {
	case models.Tracks:
		return i.fetchData(ctx, userId, models.Tracks, i.extSpotify.GetTopTracks, i.mapTrackData)
	case models.Artists:
		return i.fetchData(ctx, userId, models.Artists, i.extSpotify.GetTopArtists, i.mapArtistData)
	case models.Genres:
		// TODO: Implementar
		return nil, fmt.Errorf("genres not implemented yet")
	default:
		return nil, fmt.Errorf("invalid data type")
	}
}

// fetchData fetches the required data, either from the repo or the API.
func (i *Implementation) fetchData(
	ctx context.Context,
	userId string,
	dataType string,
	apiFunc interface{},
	mapFunc func(interface{}, string) *models.Data,
) (*models.Data, error) {
	var token *oauth2.Token
	// Intenta obtener los datos de la base de datos (sin consultar al API de Spotify)
	data, err := i.repo.GetData(userId, dataType)
	if err != nil && err != dynamodb.ErrNotFound {
		return nil, err
	}
	if data != nil && !i.dataHasExpired(data.UpdatedAt) {
		log.Info("Data is updated")
		return data, nil
	}
	// Si no hay datos en la base de datos o si los datos han expirado, intenta obtener el token del cache
	token, err = i.getTokenFromCache(ctx, userId)
	if err != nil {
		return nil, err
	}
	if token == nil {
		authUser, err := i.repo.GetAuthUser(userId)
		if err != nil {
			return nil, err
		}
		oldToken := authUser.Token
		token, err = i.renewToken(ctx, oldToken)
		if err != nil {
			return nil, err
		}
		//	TODO: guardar el token nuevo en la base de datos y en cache
		newExpiry := token.Expiry.Add(15 * time.Minute)
		duration := newExpiry.Sub(time.Now())
		go i.cacheTokens.SaveWithTTL(ctx, authUser.UserId, token, duration)
		go i.SaveAuthUser(authUser.Data, token)
	}

	// verifica que api usara y lo ejecuta
	var apiData interface{}
	switch apiFuncTyped := apiFunc.(type) {
	case func(context.Context, *oauth2.Token) ([]spotify.FullTrack, error):
		apiData, err = apiFuncTyped(ctx, token) // ejecuta la funcion para tracks
	case func(context.Context, *oauth2.Token) ([]spotify.FullArtist, error):
		apiData, err = apiFuncTyped(ctx, token) // ejecuta la funcion para artists
	}
	if err != nil {
		return nil, err
	}

	// mapea los datos obtenidos de la api
	dataToReturn := mapFunc(apiData, userId)

	// Dispara la goroutine para guardar los datos de manera asíncrona
	go func() {
		_, err := i.repo.Save(context.Background(), dataToReturn)
		if err != nil {
			// Puedes registrar el error, pero no podrás manejarlo directamente
			// en la función que llamó a fetchData.
			log.Errorf("[GET] Error al guardar los tops %s para el usuario: %s", dataType, userId)
		}
	}()

	return dataToReturn, nil
}

func (i *Implementation) dataHasExpired(updatedAt time.Time) bool {
	now := time.Now()
	ttlDays, err := strconv.Atoi(i.appConfig.RefreshUserDataTTL)
	if err != nil {
		log.Error("Error al convertir APP_REFRESH_USER_DATA_TTL a entero: %v", err)
	}
	expirationDate := updatedAt.Add(time.Duration(ttlDays) * 24 * time.Hour)
	return now.After(expirationDate)
}

func (i *Implementation) Login(ctx context.Context, callbackUrl string) (*string, string, error) {
	state, err := generateRandomState()
	if err != nil {
		return nil, "", err
	}
	callbackURL := i.getCallbackUrl(callbackUrl)
	auth := i.setupSpotifyAuth(callbackURL)
	authURL := auth.AuthURL(state)
	log.Info(authURL)

	stateKey := fmt.Sprintf(`state-%s`, state)
	stateUrlKey := fmt.Sprintf(`state-url-%s`, state)
	go i.cacheStates.Save(ctx, stateKey, state)
	go i.cacheStates.Save(ctx, stateUrlKey, callbackURL)

	return &authURL, state, nil
}

func generateRandomState() (string, error) {
	const byteSize = 32
	buf := make([]byte, byteSize)

	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(buf), nil
}

func (i *Implementation) HandleCallback(ctx context.Context, code string, state string) (*spotify.PrivateUser, error) {
	stateKey := fmt.Sprintf(`state-%s`, state)
	stateUrlKey := fmt.Sprintf(`state-url-%s`, state)
	err := i.validateState(ctx, state, stateKey)
	if err != nil {
		return nil, err
	}

	callbackURL := i.getCallbackUrlFromCache(ctx, stateUrlKey)
	auth := i.setupSpotifyAuth(callbackURL)

	token, err := auth.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	log.Info("Token: ", token)

	httpClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	client := spotify.New(httpClient)

	user, err := client.CurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	newExpiry := token.Expiry.Add(15 * time.Minute)
	duration := newExpiry.Sub(time.Now())
	go i.cacheTokens.SaveWithTTL(ctx, user.ID, token, duration)
	go i.SaveAuthUser(user, token)

	go i.cacheStates.Delete(ctx, stateKey)
	go i.cacheStates.Delete(ctx, stateUrlKey)

	log.Info("Logged in as %s\n", user.DisplayName)
	return user, nil
}

func (i *Implementation) validateState(ctx context.Context, state string, stateKey string) error {
	_, value := i.cacheStates.Get(ctx, stateKey)
	if state != value {
		return errors.New("Invalid state parameter")
	}

	return nil
}

func (i *Implementation) getCallbackUrlFromCache(ctx context.Context, stateUrlKey string) string {
	_, value := i.cacheStates.Get(ctx, stateUrlKey)
	if value == nil {
		return i.spotifyAuthConf.RedirectURI
	}
	go i.cacheStates.Delete(ctx, stateUrlKey)
	return value.(string)
}

func (i *Implementation) getCallbackUrl(callbackUrl string) string {
	if callbackUrl == "" {
		callbackUrl = i.spotifyAuthConf.RedirectURI
	}
	redirectUrl, err := url.Parse(callbackUrl)
	if err != nil {
		log.Error(err)
	}
	q := redirectUrl.Query()
	q.Set("provider", "spotify")
	redirectUrl.RawQuery = q.Encode()
	return redirectUrl.String()
}

func (i *Implementation) setupSpotifyAuth(callbackUrl string) *spotifyauth.Authenticator {
	spotifyAuthConf := i.spotifyAuthConf.GetSpotifyAuthConf()

	scopes := []string{
		spotifyauth.ScopeUserTopRead,
		spotifyauth.ScopeUserReadPrivate,
		spotifyauth.ScopeUserReadEmail,
		spotifyauth.ScopeUserReadRecentlyPlayed,
	}

	authenticator := spotifyauth.New(
		spotifyauth.WithRedirectURL(callbackUrl),
		spotifyauth.WithScopes(scopes...),
		spotifyauth.WithClientID(spotifyAuthConf.ClientID),
		spotifyauth.WithClientSecret(spotifyAuthConf.ClientSecret),
	)
	return authenticator
}

func (i *Implementation) SaveAuthUser(data *spotify.PrivateUser, token *oauth2.Token) {
	authUser := i.mapAuthUserData(data, token)
	_, err := i.repo.SaveAuthUser(context.Background(), authUser)
	if err != nil {
		log.Errorf("[GET] Error al guardar los datos del usuario: %s", data.ID)
	}
}

func (i *Implementation) mapAuthUserData(data *spotify.PrivateUser, token *oauth2.Token) *models.AuthUserData {
	return &models.AuthUserData{
		UserId: data.ID,
		Token:  token,
		Data:   data,
	}
}

func (i *Implementation) mapTrackData(data interface{}, userId string) *models.Data {
	tracks, ok := data.([]spotify.FullTrack)
	if !ok {
		return nil
	}
	tDetails := make([]models.DataDetails, len(tracks))
	for index, track := range tracks {
		artists := make([]models.Artist, len(track.Artists))
		for j, artist := range track.Artists {
			artists[j] = models.Artist{
				Name: artist.Name,
				ID:   artist.ID.String(),
			}
		}

		tDetails[index] = models.TrackDetails{
			ID:          track.ID.String(),
			Name:        track.Name,
			Artists:     artists,
			ReleaseDate: track.Album.ReleaseDate,
		}
	}

	return &models.Data{
		UserId:   userId,
		DataType: models.Tracks,
		Data:     tDetails,
		Source:   Spotify,
		Count:    len(tracks),
	}
}

func (i *Implementation) mapArtistData(data interface{}, userId string) *models.Data {
	artists, ok := data.([]spotify.FullArtist)
	if !ok {
		return nil
	}
	aDetails := make([]models.DataDetails, len(artists))
	for index, artist := range artists {
		aDetails[index] = models.ArtistDetails{
			ID:       artist.ID.String(),
			Name:     artist.Name,
			Genres:   artist.Genres,
			Realname: artist.Name,
		}
	}

	return &models.Data{
		UserId:   userId,
		DataType: models.Artists,
		Data:     aDetails,
		Source:   Spotify,
		Count:    len(artists),
	}
}

func (i *Implementation) getTokenFromCache(ctx context.Context, userId string) (*oauth2.Token, error) {
	/*
		//si estamos en local devolver un mock
		if util.IsLocal(util.GetEnvironment()) {
			return &oauth2.Token{
				//nolint:lll
				AccessToken: "BQAKWQias6Pf1KFNb02kjcQkpxxdO9iRcS1vNf2zPnXUP1wM4xv8fexFLBtJxJxMS0CVcZ0UNF3mRpvKlkhtqKYmXdMVt6ERyRphAWp8mL9qMQBFYMBOfghzVpXn4Ld4URW9ZEtxT8a0S0eYojLHGoDBjqMe6XqUFTIsfgnHNmRm0bjC0gDU4IE3zTvHKbQuHq610H_LY5qWW3-IzavoliZK8A",
				TokenType:   "Bearer",
			}, nil
		}
	*/
	_, value := i.cacheTokens.Get(ctx, userId)

	if value == nil {
		fmt.Printf("token not found in cache for user: %s", userId)
		return nil, nil
	}

	validToken, ok := value.(*oauth2.Token)
	if !ok {
		return nil, fmt.Errorf("cached value is not a valid oauth2 token for user: %s", userId)
	}
	return validToken, nil
}

func (i *Implementation) renewToken(ctx context.Context, token *oauth2.Token) (*oauth2.Token, error) {
	conf := &oauth2.Config{}
	src := conf.TokenSource(ctx, token)
	newToken, err := src.Token()
	if err != nil {
		log.Error("Error renovando el token: %v", err)
		return nil, fmt.Errorf("No se pudo renovar el token")
	}
	return newToken, nil
}
