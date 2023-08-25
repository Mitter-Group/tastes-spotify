package entity

// Config for application
type Config struct {
	Env          string      `json:"-"`
	Version      string      `json:"string"`
	OAuthExample OAuthConfig `json:"oauth_example"`
}

type OAuthConfig struct {
	Audience        string `json:"audience"`
	Domain          string `json:"domain"`
	ClientID        string `json:"client_id"`
	ClientSecretKey string `json:"client_secret_key"`
	TokenURL        string `json:"token_url"`
	Scope           string `json:"scope"`
}
