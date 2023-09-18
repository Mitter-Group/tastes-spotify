package external

import (
	"context"
	"log"

	"github.com/chunnior/spotify/internal/models"
	"golang.org/x/oauth2"

	"github.com/zmb3/spotify/v2"
)

type Implementation struct {
	spotifyClient *spotify.Client
}

func NewIntegration(cfg models.Config) Integration {
	ctx := context.Background()

	token, err := cfg.Spotify.GetSpotifyAuthConf().Token(ctx)
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	client := initializeSpotifyClient(ctx, token)

	return &Implementation{
		spotifyClient: client,
	}
}

func initializeSpotifyClient(ctx context.Context, token *oauth2.Token) *spotify.Client {
	httpClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	return spotify.New(httpClient)
}
func (i *Implementation) GetTopTracks(ctx context.Context, token *oauth2.Token) ([]spotify.FullTrack, error) {
	client := initializeSpotifyClient(ctx, token)
	options := []spotify.RequestOption{
		spotify.Timerange(spotify.LongTermRange),
	}

	topTracks, err := client.CurrentUsersTopTracks(ctx, options...)
	if err != nil {
		return nil, err
	}
	return topTracks.Tracks, nil
}

func (i *Implementation) GetTopArtists(ctx context.Context, token *oauth2.Token) ([]spotify.FullArtist, error) {
	client := initializeSpotifyClient(ctx, token)
	options := []spotify.RequestOption{
		spotify.Timerange(spotify.LongTermRange),
	}

	topArtists, err := client.CurrentUsersTopArtists(ctx, options...)
	if err != nil {
		return nil, err
	}
	return topArtists.Artists, nil
}

func (i *Implementation) GetUserDetails(ctx context.Context) (*spotify.PrivateUser, error) {
	user, err := i.spotifyClient.CurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}
