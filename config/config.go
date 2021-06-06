package config

import (
	"arcdps/logger"
	"encoding/json"
	"os"
)

type Config struct {
	Destination     string `json:"destination" db:"destination"`
	URL             string `json:"url" db:"url"`
	Filename        string `json:"filename" db:"filename"`
	Gw2LauncherPath string `json:"gw2LauncherPath" db:"gw2LauncherPath"`
}

func ReadConfig(filepath string) (*Config, error) {
	log := logger.Logger()
	file, err := os.Open(filepath)
	if err != nil {
		log.Error().Msg("Config file not present, please create a config.json")
		return nil, err
	}

	defer file.Close()
	decoder := json.NewDecoder(file)
	config := Config{}

	err = decoder.Decode(&config)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	return &config, nil
}
