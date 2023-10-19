package spotify

import (
	"context"
	"time"

	"github.com/chunnior/spotify/internal/models"
	dynamodb "github.com/chunnior/spotify/pkg/aws/dynamodb"
)

type DataRepo struct {
	dynamo               dynamodb.Client
	cfg                  models.AWS
	tableUserSpotify     string
	tableSpotifyAuthUser string
}

func NewConfig(dynamoTable dynamodb.Client, cfg models.Config) Spec {
	return &DataRepo{
		dynamo:               dynamoTable,
		cfg:                  *cfg.AWS,
		tableUserSpotify:     cfg.AWS.UserSpotifyData.TableName,
		tableSpotifyAuthUser: cfg.AWS.SpotifyAuthUser.TableName,
	}
}

func (dr *DataRepo) Save(ctx context.Context, dataReq *models.Data) (models.Data, error) {
	now := time.Now()

	// Si `CreatedAt` no ha sido establecido (es decir, es un nuevo objeto), entonces se asigna la fecha y hora actual.
	if dataReq.CreatedAt.IsZero() {
		dataReq.CreatedAt = now
	}

	// Independientemente de si es un nuevo objeto o una actualización, `UpdatedAt` siempre se establecerá a la hora actual.
	dataReq.UpdatedAt = now

	err := dr.dynamo.Save(dr.tableUserSpotify, dataReq)
	if err != nil {
		return models.Data{}, err
	}

	return *dataReq, nil
}

func (dr *DataRepo) SaveAuthUser(ctx context.Context, dataReq *models.AuthUserData) (models.AuthUserData, error) {
	now := time.Now()

	// Si `CreatedAt` no ha sido establecido (es decir, es un nuevo objeto), entonces se asigna la fecha y hora actual.
	if dataReq.CreatedAt.IsZero() {
		dataReq.CreatedAt = now
	}

	// Independientemente de si es un nuevo objeto o una actualización, `UpdatedAt` siempre se establecerá a la hora actual.
	dataReq.UpdatedAt = now

	err := dr.dynamo.Save(dr.tableUserSpotify, dataReq)
	if err != nil {
		return models.AuthUserData{}, err
	}

	return *dataReq, nil
}

func (dr *DataRepo) GetData(userId string, dataType string) (*models.Data, error) {

	var data *models.Data

	err := dr.dynamo.GetOneWithSort(dr.tableUserSpotify, userId, dataType, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
