package config

import (
	"encoding/json"
	"os"
)

type RepositoryConfig struct {
	Type   string          `json:"type"`
	Config json.RawMessage `json:"config"`
}

type Config struct {
	SessionRepository RepositoryConfig `json:"sessionRepository"`
	UserRepository    RepositoryConfig `json:"userRepository"`
}

func ConfigFromJson(filePath string) (Config, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
