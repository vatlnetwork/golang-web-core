package config

import (
	"encoding/json"
	"os"
)

type Secrets struct {
	MongoDatabaseUsername string `json:"mongoDatabaseUsername"`
	MongoDatabasePassword string `json:"mongoDatabasePassword"`
}

func SecretsFromJson(filePath string) (Secrets, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return Secrets{}, err
	}

	var secrets Secrets
	err = json.Unmarshal(bytes, &secrets)
	if err != nil {
		return Secrets{}, err
	}

	return secrets, nil
}
