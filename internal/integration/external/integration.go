package external

import (
	"context"
	"log"

	"github.com/chunnior/spotify/internal/models"
	spotifyAuth "github.com/zmb3/spotify/v2/auth"

	"github.com/zmb3/spotify/v2"
)

type Implementation struct {
	spotifyClient *spotify.Client
}

func NewIntegration(cfg models.Config) Integration {

	ctx := context.Background()

	configCred := cfg.Spotify.GetSpotifyAuthConf()
	token, err := configCred.Token(ctx)
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	httpClient := spotifyAuth.New().Client(ctx, token)
	client := spotify.New(httpClient)

	return &Implementation{
		spotifyClient: client,
	}
}

func (i *Implementation) GetTopTracks(ctx context.Context) ([]spotify.FullTrack, error) {

	options := []spotify.RequestOption{
		spotify.Timerange(spotify.LongTermRange),
	}

	topTracks, err := i.spotifyClient.CurrentUsersTopTracks(ctx, options...)
	if err != nil {
		return nil, err
	}
	return topTracks.Tracks, nil
}
