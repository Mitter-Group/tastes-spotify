package spotify

import (
	"context"
	"crypto/rand"
	"math/big"

	"github.com/chunnior/spotify/internal/integration/external"
	"github.com/chunnior/spotify/internal/models"
	configRepo "github.com/chunnior/spotify/internal/repository/data"
	"github.com/chunnior/spotify/internal/util/log"
	"github.com/chunnior/spotify/pkg/aws/dynamodb"
)

type Implementation struct {
	repo       configRepo.Spec
	extSpotify external.Integration
}

func NewUseCase(repo configRepo.Spec, external external.Integration) *Implementation {
	return &Implementation{
		repo:       repo,
		extSpotify: external,
	}
}

func (i *Implementation) Save(ctx context.Context, tasteReq models.DataRequest) (models.Data, error) {

	randInt, _ := rand.Int(rand.Reader, big.NewInt(100))

	taste := models.Data{
		UserId: tasteReq.UserId,
		Data: []models.DataDetails{
			{
				ID:   randInt.String(),
				Name: "Track 1",
			},
		},
	}

	return taste, nil
}

func (i *Implementation) Get(ctx context.Context, dataType string, userId string) (*models.Data, error) {

	switch dataType {
	case "tracks":
		log.InfofWithContext(ctx, "[GET] Getting top tracks for user: %s", userId)
		return i.getTopTracks(ctx, userId)
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

func (i *Implementation) getTopTracks(ctx context.Context, userId string) (*models.Data, error) {

	// TODO: Implementar llamada a API de Spotify
	// Debemos buscar la data en memoria o base de datos, si no existe, llamar a la API de Spotify

	data, err := i.repo.GetData(userId, "tracks")
	if err != nil {
		if err == dynamodb.ErrNotFound {
			log.Errorf("[GET] Data not found for user: %s", userId)

			tracks, err := i.extSpotify.GetTopTracks(ctx)

			if err != nil {
				log.Errorf("[GET] Error getting top tracks for user: %s", userId)
				return nil, err
			}

			log.Info(tracks)

			return nil, nil
		}
		return nil, err
	}

	if data != nil {
		return data, nil
	}

	return data, nil
}
