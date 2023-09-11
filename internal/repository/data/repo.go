package spotify

import (
	"context"

	"github.com/chunnior/spotify/internal/models"
	dynamodb "github.com/chunnior/spotify/pkg/aws/dynamodb"
)

type DataRepo struct {
	dynamo           dynamodb.Client
	cfg              models.AWS
	tableUserSpotify string
}

func NewConfig(dynamoTable dynamodb.Client, cfg models.Config) Spec {
	return &DataRepo{
		dynamo:           dynamoTable,
		cfg:              *cfg.AWS,
		tableUserSpotify: cfg.AWS.UserSpotifyData.TableName,
	}
}

func (dr *DataRepo) Save(ctx context.Context, tasteReq models.DataRequest) (models.Data, error) {
	return models.Data{}, nil
}

func (dr *DataRepo) GetData(userId string, dataType string) (*models.Data, error) {

	var data *models.Data

	err := dr.dynamo.GetOneWithSort(dr.tableUserSpotify, userId, dataType, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
