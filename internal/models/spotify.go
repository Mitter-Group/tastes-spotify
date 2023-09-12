package models

import (
	"os"

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
