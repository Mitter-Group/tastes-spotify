package spotify

import (
	"context"

	"github.com/chunnior/spotify/internal/models"
)

type Spec interface {
	GetData(userId string, dataType string) (*models.Data, error)
	GetAuthUser(userId string) (*models.AuthUserData, error)
	Save(ctx context.Context, dataReq *models.Data) (models.Data, error)
	SaveAuthUser(ctx context.Context, dataReq *models.AuthUserData) (models.AuthUserData, error)
}
