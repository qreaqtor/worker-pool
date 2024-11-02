package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	configPath       = "CONFIG_PATH"
	defultConfigPath = "./config/config.yaml"
)

// Loads the configuration path from the CONFIG_PATH env.
// If the path is empty or does not exist, the default configuration path is used ("./config/config.yaml").
func Load() (*Config, error) {
	path, exists := os.LookupEnv(configPath)
	if !exists || path == "" {
		path = defultConfigPath
	}

	return LoadFromPath(path)
}

// Loads the configuration from provided path.
func LoadFromPath(path string) (*Config, error) {
	cfg := &Config{}

	_, err := os.Stat(path)
	if err != nil {
		return nil, errNotExists
	}

	err = cleanenv.ReadConfig(path, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
