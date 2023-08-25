package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type Response struct {
	Message string `json:"message"`
}

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

var jwks = Jwks{}
var ScopesRequired = []string{}
var AudiencesAllowed = []string{}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

func validationKeyGetter(token *jwt.Token) (interface{}, error) {
	// Verify all audience
	validAudience := false
	for _, v := range AudiencesAllowed {
		checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(v, true)
		if checkAud {
			validAudience = true
			break
		}
	}
	if !validAudience {
		return token, errors.New("invalid audience")
	}

	// Verify 'iss' claim
	iss := os.Getenv("AUTH_ISS")
	checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
	if !checkIss {
		return token, errors.New("invalid issuer")
	}

	cert, err := getPemCert(token)
	if err != nil {
		return nil, err
	}

	result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
	return result, nil
}

func Initialize() {
	resp, err := http.Get(os.Getenv("AUTH_ISS") + ".well-known/jwks.json")

	if err != nil {
		log.Panic("Fatal initializing oauth library: ", err.Error())
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		log.Panic("Fatal initializing oauth library: ", err.Error())
	}
	osAudiencesAllowed := os.Getenv("AUTH_AUDIENCE")
	if len(osAudiencesAllowed) < 1 {
		log.Panic("Fatal initializing oauth library: audiences not configured")
	}
	AudiencesAllowed = strings.Split(osAudiencesAllowed, ",")

	osScopeRequired := os.Getenv("AUTH_SCOPE_REQUIRED")
	if len(osScopeRequired) > 0 {
		ScopesRequired = strings.Split(osScopeRequired, " ")
	}
}

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""

	for k := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("unable to find appropriate key")
		return cert, err
	}

	return cert, nil
}

func extractor(c *fiber.Ctx) (string, error) {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return "", nil
	}
	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", errors.New("Authorization header format must be Bearer {token}")
	}
	return authHeaderParts[1], nil
}

func checkJWT(c *fiber.Ctx) bool {
	// Exctract token from header
	token, err := extractor(c)
	if err != nil || token == "" {
		fmt.Printf("Error extracting token: %v", err)
		return false
	}

	// Parse the token
	parsedToken, err := jwt.Parse(token, validationKeyGetter)
	if err != nil {
		fmt.Printf("Error parsing token: %v", err)
		return false
	}

	// Validate token
	if !parsedToken.Valid {
		fmt.Printf("Invalid token.")
		return false
	}

	// Get claims
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Printf("Cannot get claims from token: %v", err)
		return false
	}

	// Validate required scopes
	if !hasValidScopes(claims) {
		return false
	}

	// Adding scopes and client_id to context
	c.Locals("scopes", fmt.Sprintf("%v", claims["scope"]))
	c.Locals("client_id", fmt.Sprintf("%v", claims["client_id"]))

	return parsedToken.Valid
}

func hasValidScopes(claims jwt.MapClaims) bool {
	if len(ScopesRequired) < 1 {
		return true
	}

	scopes := fmt.Sprintf("%v", claims["scope"])
	scopesSplitted := strings.Split(scopes, " ")

	if len(scopesSplitted) < 1 {
		return false
	}

	for _, v := range ScopesRequired {
		if !contains(scopesSplitted, v) {
			fmt.Printf("Scope %v not found.", v)
			return false
		}
	}

	return true
}

// Protected does check your JWT token and validates it
func Protected(c *fiber.Ctx) error {
	if checkJWT(c) {
		return c.Next()
	}
	return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized.")
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
