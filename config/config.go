package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Destination string `json:"destination,omitempty" db:"destination"`
	URL         string `json:"url,omitempty" db:"url"`
}

func ReadConfig(filepath string) (*Config, error) {
	file, err := os.Open(filepath)
	defer file.Close()
	if err != nil {
		log.Println("Config file not present, please create a config.json")
		return nil, err
	}
	decoder := json.NewDecoder(file)
	config := Config{}

	err = decoder.Decode(&config)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &config, nil
}
