package config

import (
	"arcdps/logger"
	"encoding/json"
	"os"
)

type DxOptions struct {
	FileName    string `json:"fileName,omitempty" db:"fileName"`
	Destination string `json:"destination,omitempty" db:"destination"`
}
type Config struct {
	Destination       string    `json:"destination,omitempty" db:"destination"`
	URL               string    `json:"url,omitempty" db:"url"`
	Filename          string    `json:"filename,omitempty" db:"filename"`
	Gw2LauncherPath   string    `json:"gw2LauncherPath,omitempty" db:"gw2_launcher_path"`
	EnableGw2Launcher bool      `json:"enableGw2Launcher,omitempty" db:"enable_gw2_launcher"`
	LogLevel          string    `json:"logLevel,omitempty" db:"log_level"`
	Dx11              DxOptions `json:"dx11,omitempty" db:"dx11"`
	Dx9               DxOptions `json:"dx9,omitempty" db:"dx9"`
	EnableDx11        bool      `json:"enableDx11,omitempty" db:"enableDx11"`
	RetainOldVersion  bool      `json:"retainOldVersion,omitempty" db:"retainOldVersion"`
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
