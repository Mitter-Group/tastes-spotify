package spotify

import (
	"context"
	"crypto/rand"
	"math/big"

	"github.com/chunnior/spotify/internal/models"
	"github.com/chunnior/spotify/internal/util/log"
)

type Implementation struct {
}

func NewUseCase() *Implementation {
	return &Implementation{}
}

func (h *Implementation) Save(ctx context.Context, tasteReq models.TasteRequest) (models.Taste, error) {

	randInt, _ := rand.Int(rand.Reader, big.NewInt(100))

	taste := models.Taste{
		ID:   randInt.String(),
		Name: tasteReq.Name,
	}

	return taste, nil
}

func (h *Implementation) Get(ctx context.Context, dataType string, userId string) (*models.DataResponse, error) {

	switch dataType {
	case "tracks":
		log.InfofWithContext(ctx, "[GET] Getting top tracks for user: %s", userId)
		return h.getTopTracks(userId)
	case "artists":
		log.InfofWithContext(ctx, "[GET] Getting top artists for user: %s", userId)
		return nil, nil
	case "genres":
		log.InfofWithContext(ctx, "[GET] Getting top genres for user: %s", userId)
		return nil, nil
	default:
		return nil, nil
	}

}

func (h *Implementation) getTopTracks(userId string) (*models.DataResponse, error) {

	// TODO: Implementar llamada a API de Spotify
	// Debemos buscar la data en memoria o base de datos, si no existe, llamar a la API de Spotify

	data := models.DataResponse{
		UserId: userId,
		Data: []models.Taste{
			{
				ID:   "1",
				Name: "Track 1",
			},
			{
				ID:   "2",
				Name: "Track 2",
			},
		},
		Source: "SPOTIFY",
	}

	return &data, nil
}
