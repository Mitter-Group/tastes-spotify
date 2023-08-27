package tastes

import (
	"context"
	"crypto/rand"
	"math/big"

	"github.com/chunnior/spotify/internal/models"
)

type Implementation struct {
}

func NewUseCase() *Implementation {
	return &Implementation{}
}

// Save is a function that saves a taste
func (h *Implementation) Save(ctx context.Context, tasteReq models.TasteRequest) (models.Taste, error) {

	// generate random id
	randInt, _ := rand.Int(rand.Reader, big.NewInt(100))

	taste := models.Taste{
		ID:   randInt.String(),
		Name: tasteReq.Name,
	}

	return taste, nil
}
