package tastes

import (
	"context"

	"github.com/chunnior/spotify/internal/models"
)

type UseCase interface {
	Save(context.Context, models.TasteRequest) (models.Taste, error)
}
