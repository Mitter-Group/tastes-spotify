package util

import (
	"fmt"
	"os"
	"strings"
)

var environment string = ""

func GetEnvironment() string {
	if environment != "" {
		return environment
	}

	osValue := os.Getenv("ENV")
	if osValue != "" {
		environment = strings.ToLower(osValue)
	} else {
		environment = "local"
	}

	fmt.Println("Running environment: ", environment)
	return environment
}

func IsDevelopmentEnvironment() bool {
	currentEnvironment := GetEnvironment()
	return IsDevelopmentOrLocal(currentEnvironment)
}

func IsProductionEnvironment() bool {
	return !IsDevelopmentEnvironment()
}

func IsLocal(env string) bool {
	env = strings.ToLower(env)
	return env == "local"
}
func IsDevelopment(env string) bool {
	env = strings.ToLower(env)
	return env == "dev" || env == "devcde" || env == "development"
}
func IsStage(env string) bool {
	return env == "stage" || env == "stg" || env == "stgcde"
}
func IsProduction(env string) bool {
	return env == "prd" || env == "prod" || env == "production" || env == "cde"
}
func IsDevelopmentOrLocal(env string) bool {
	env = strings.ToLower(env)
	return IsLocal(env) || IsDevelopment(env)
}
