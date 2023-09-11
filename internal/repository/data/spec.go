package spotify

import (
	"context"

	"github.com/chunnior/spotify/internal/models"
)

type Spec interface {
	GetData(userId string, dataType string) (*models.Data, error)
	Save(ctx context.Context, tasteReq models.DataRequest) (models.Data, error)
}
