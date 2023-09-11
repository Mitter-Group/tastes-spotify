package oauth

import (
	"context"
	"net/url"
	"os"

	"github.com/chunnior/spotify/internal/models"
	"github.com/chunnior/spotify/internal/util/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type OAuthClient struct {
	clientID        string
	tokenURL        string
	clientSecretKey string
	audience        string
	scope           string
	tokenSource     oauth2.TokenSource
}

func NewOAuthClient(cfg models.OAuthConfig) *OAuthClient {
	return &OAuthClient{
		clientID:        cfg.ClientID,
		tokenURL:        cfg.TokenURL,
		clientSecretKey: cfg.ClientSecretKey,
		audience:        cfg.Audience,
		scope:           cfg.Scope,
		tokenSource:     nil,
	}
}

func (o *OAuthClient) GetAccessToken() (string, error) {

	if o.tokenSource != nil {
		return GetAccessTokenFromSource(o.tokenSource)
	}

	log.Debug("generating token source for first time")
	clientSecret := os.Getenv(o.clientSecretKey)

	scopes := []string{}
	if o.scope != "" {
		scopes = append(scopes, o.scope)
	}

	oauthConfig := clientcredentials.Config{
		ClientID:       o.clientID,
		ClientSecret:   clientSecret,
		Scopes:         scopes,
		TokenURL:       o.tokenURL,
		EndpointParams: url.Values{"audience": {o.audience}},
	}

	tokenSorce := oauthConfig.TokenSource(context.Background())
	o.tokenSource = tokenSorce

	return GetAccessTokenFromSource(o.tokenSource)
}

func GetAccessTokenFromSource(tokenSorce oauth2.TokenSource) (string, error) {
	token, err := tokenSorce.Token()

	if err != nil {
		log.Error("error on get access token", err)
		return "", err
	}

	return token.AccessToken, err
}
