package models

import (
	configAws "github.com/chunnior/spotify/pkg/aws"
	dynamodb "github.com/chunnior/spotify/pkg/aws/dynamodb"
)

// Config for application
type Config struct {
	Env          string          `json:"-"`
	Version      string          `json:"string"`
	OAuthExample OAuthConfig     `json:"oauth_example"`
	AWS          *AWS            `json:"aws"`
	Spotify      SpotifyAuthConf `json:"spotify"`
}

type OAuthConfig struct {
	Audience        string `json:"audience"`
	Domain          string `json:"domain"`
	ClientID        string `json:"client_id"`
	ClientSecretKey string `json:"client_secret_key"`
	TokenURL        string `json:"token_url"`
	Scope           string `json:"scope"`
}

type AWS struct {
	Credentials     configAws.AWSCredentials `json:"credentials"`
	UserSpotifyData dynamodb.DynamoTable     `json:"user_spotify_data"`
	LocalstackPort  string                   `json:"localstack_port"`
}
