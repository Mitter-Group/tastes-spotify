package models

import (
	"fmt"
	"os"
	"time"

	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
)

type SpotifyAuthConf struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	TokenURL     string `json:"token_url"`
	RedirectURI  string `json:"redirect_uri"`
}

func (s SpotifyAuthConf) GetSpotifyAuthConf() *clientcredentials.Config {
	return &clientcredentials.Config{
		ClientID:     os.Getenv(s.ClientID),
		ClientSecret: os.Getenv(s.ClientSecret),
		TokenURL:     spotifyauth.TokenURL,
	}
}

const (
	Tracks  string = "tracks"
	Artists string = "artists"
	Genres  string = "genres"
)

type DataRequest struct {
	DataType string `json:"data_type" example:"tracks"`
	UserId   string `json:"user_id" example:"123456789"`
}

type DynamoDBMetadata struct {
	CreatedAt time.Time `dynamodbav:"createdAt"`
	UpdatedAt time.Time `dynamodbav:"updatedAt"`
}

type Data struct {
	UserId           string        `json:"user_id" dynamodbav:"userID"`
	DataType         string        `json:"data_type" dynamodbav:"dataType"`
	Data             []DataDetails `json:"data"`
	Source           string        `json:"source" example:"SPOTIFY"`
	Count            int           `json:"count"`
	DynamoDBMetadata `dynamodbav:",inline"`
}

type AuthUserData struct {
	UserId           string      `json:"user_id" dynamodbav:"userID"`
	SpotifyUserId    string      `json:"spotify_user_id" dynamodbav:"spotifyUserID"`
	TokenType        string      `json:"token_type" example:"bearer"`
	RefreshToken     string      `json:"refresh_token"`
	TokenExpiration  time.Time   `json:"token_expiration"`
	Data             DataDetails `json:"data"`
	DynamoDBMetadata `dynamodbav:",inline"`
}

type DataDetails interface{}

type TrackDetails struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Artists     []Artist `json:"artist"`
	ReleaseDate string   `json:"release_date"`
}

type Artist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ArtistDetails struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Genres   []string `json:"genre"`
	Realname string   `json:"realname"`
}

type GenreDetails struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (dr *DataRequest) Validate() error {
	switch dr.DataType {
	case Tracks, Artists, Genres:
		return nil
	default:
		return fmt.Errorf("Invalid DataType value: %s", dr.DataType)
	}
}
