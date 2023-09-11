package util

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/chunnior/spotify/internal/models"
	"github.com/chunnior/spotify/internal/util/log"
)

func ReadConfig(env string) (*models.Config, error) {
	cfg := models.Config{}

	configFile := readConfigFile(env)
	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(&cfg); err != nil {
		log.Error("parsing config file", err.Error())
		return nil, err
	}

	cfg.Env = env

	// Add this cfg to environment
	os.Setenv("AUTH0_AUDIENCE", cfg.OAuthExample.Audience)
	os.Setenv("AUTH0_DOMAIN", cfg.OAuthExample.Domain)

	return &cfg, nil
}

func readConfigFile(env string) *os.File {

	// using the function
	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	fileNames := []string{
		mydir + "/config/" + env + ".json",
	}

	for _, name := range fileNames {
		configFile, err := os.Open(name)
		if err != nil {
			fmt.Println("opening config file error: ", err.Error())
		} else {
			fmt.Println("Config file", name, " loaded")
			return configFile
		}
	}

	return nil
}
