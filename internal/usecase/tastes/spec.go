package tastes

import (
	"context"

	"github.com/chunnior/tastes/internal/models"
)

type UseCase interface {
	Save(context.Context, models.TasteRequest) (models.Taste, error)
}
