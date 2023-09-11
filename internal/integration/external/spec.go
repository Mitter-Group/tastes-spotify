package external

import (
	"context"

	"github.com/zmb3/spotify/v2"
)

type Integration interface {
	GetTopTracks(ctx context.Context) ([]spotify.FullTrack, error)
}
