package spotify

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/chunnior/spotify/internal/integration/external"
	"github.com/chunnior/spotify/internal/models"
	configRepo "github.com/chunnior/spotify/internal/repository/data"
	"github.com/chunnior/spotify/internal/util"
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
	cacheTokens     cache.Spec
	cacheStates     cache.Spec
}

const (
	Spotify     = "SPOTIFY"
	MemoryCache = "tokens"
	CacheSize   = 2000
	CacheExpiry = 3540
)

func NewUseCase(spotifyAuthConf models.SpotifyAuthConf, repo configRepo.Spec, external external.Integration) *Implementation {
	return &Implementation{
		repo:            repo,
		extSpotify:      external,
		spotifyAuthConf: spotifyAuthConf,
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

	// intenta obtener los datos de la base de datos
	data, err := i.repo.GetData(userId, dataType)
	if err != nil && err != dynamodb.ErrNotFound {
		return nil, err
	}

	// TODO: Verificar si los datos estan actualizados

	if data != nil {
		return data, nil
	}

	// obtiene el token del cache
	token, err := i.getTokenFromCache(ctx, userId)
	if err != nil {
		return nil, err
	}

	// TODO: si el token no existe en el cache, buscar uno nuevo con el refresh token

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

func (i *Implementation) Login(ctx context.Context) (*string, string, error) {
	state, err := generateRandomState()
	if err != nil {
		return nil, "", err
	}

	auth := i.setupSpotifyAuth()

	authURL := auth.AuthURL(state)
	log.Info(authURL)

	stateKey := fmt.Sprintf(`state-%s`, state)
	go i.cacheStates.Save(ctx, stateKey, state)

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
	err := i.validateState(ctx, state)
	if err != nil {
		return nil, err
	}

	auth := i.setupSpotifyAuth()

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

	go i.cacheTokens.Save(ctx, user.ID, token)

	log.Info("Logged in as %s\n", user.DisplayName)
	return user, nil
}

func (i *Implementation) validateState(ctx context.Context, state string) error {
	stateKey := fmt.Sprintf(`state-%s`, state)
	_, value := i.cacheStates.Get(ctx, stateKey)
	if state != value {
		return errors.New("Invalid state parameter")
	}
	return nil
}

func (i *Implementation) setupSpotifyAuth() *spotifyauth.Authenticator {
	spotifyAuthConf := i.spotifyAuthConf.GetSpotifyAuthConf()

	scopes := []string{
		spotifyauth.ScopeUserTopRead,
		spotifyauth.ScopeUserReadPrivate,
		spotifyauth.ScopeUserReadEmail,
		spotifyauth.ScopeUserReadRecentlyPlayed,
	}

	return spotifyauth.New(
		spotifyauth.WithRedirectURL(i.spotifyAuthConf.RedirectURI),
		spotifyauth.WithScopes(scopes...),
		spotifyauth.WithClientID(spotifyAuthConf.ClientID),
		spotifyauth.WithClientSecret(spotifyAuthConf.ClientSecret),
	)
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

	//si estamos en local devolver un mock
	if util.IsLocal(util.GetEnvironment()) {
		return &oauth2.Token{
			//nolint:lll
			AccessToken: "BQCvOZkUP7K77vMU7YgaSMJOD5Pr2_omud-_gdeC1VzNUqkA-iWYQEBHg-eDJz1fE0Qe2XgVb6lkUK44NshXq_uBa-F42wbKxHI_qHJNwAVWuWdGDP0yXzGvXc8BPoWcluxyHO_yEjL0EsiiCmYTy108PZT6XfsrxGbBtBLBSTIN43czvnLHPv8qK2KIjTa8EvZu1GqNcwKiQ9qYxYkb2oHnORXKDrRK4A",
			TokenType:   "Bearer",
		}, nil
	}
	_, value := i.cacheTokens.Get(ctx, userId)

	if value == nil {
		return nil, fmt.Errorf("token not found in cache or expired for user: %s", userId)
	}

	validToken, ok := value.(*oauth2.Token)
	if !ok {
		return nil, fmt.Errorf("cached value is not a valid oauth2 token for user: %s", userId)
	}

	return validToken, nil
}
