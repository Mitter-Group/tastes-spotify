package external

import (
	"context"

	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
)

type Integration interface {
	GetTopTracks(ctx context.Context, token *oauth2.Token) ([]spotify.FullTrack, error)
	GetTopArtists(ctx context.Context, token *oauth2.Token) ([]spotify.FullArtist, error)
	GetUserDetails(ctx context.Context) (*spotify.PrivateUser, error)
}
