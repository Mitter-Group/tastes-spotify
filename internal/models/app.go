package models

import "os"

type AppConfig struct {
	RefreshUserDataTTL string `json:"refresh_user_data_ttl"`
}

func (a AppConfig) GetSpotifyAuthConf() *AppConfig {
	return &AppConfig{
		RefreshUserDataTTL: os.Getenv(a.RefreshUserDataTTL),
	}
}
