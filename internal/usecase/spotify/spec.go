package spotify

import (
	"context"

	"github.com/chunnior/spotify/internal/models"
	"github.com/zmb3/spotify/v2"
)

type UseCase interface {
	Save(ctx context.Context, tasteReq models.DataRequest) (models.Data, error)
	Get(ctx context.Context, dataType string, userId string) (*models.Data, error)
	Login(ctx context.Context) (*string, string, error)
	HandleCallback(ctx context.Context, code string, state string) (*spotify.PrivateUser, error)
}
