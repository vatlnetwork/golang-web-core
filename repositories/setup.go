package repositories

import (
	"encoding/json"
	"fmt"
	"golang-web-core/config"
	"golang-web-core/domain"
	"golang-web-core/logging"
	"golang-web-core/repositories/sessionrepo"
	"golang-web-core/repositories/userrepo"
)

type Repositories struct {
	SessionRepository domain.SessionRepository
	UserRepository    domain.UserRepository
}

func SetupRepositories(config config.Config, logger *logging.Logger) (Repositories, error) {
	repositories := Repositories{}

	var sessionRepository domain.SessionRepository
	switch config.SessionRepository.Type {
	case "MongoSessionRepository":
		var mongoSessionRepoConfig sessionrepo.MongoSessionRepoConfig
		err := json.Unmarshal(config.SessionRepository.Config, &mongoSessionRepoConfig)
		if err != nil {
			return Repositories{}, err
		}
		mongoSessionRepository, err := sessionrepo.NewMongoSessionRepository(mongoSessionRepoConfig, logger)
		if err != nil {
			return Repositories{}, err
		}
		sessionRepository = &mongoSessionRepository
	default:
		return Repositories{}, fmt.Errorf("invalid session repository type: %s", config.SessionRepository.Type)
	}
	repositories.SessionRepository = sessionRepository

	var userRepository domain.UserRepository
	switch config.UserRepository.Type {
	case "MongoUserRepository":
		var mongoUserRepoConfig userrepo.MongoUserRepoConfig
		err := json.Unmarshal(config.UserRepository.Config, &mongoUserRepoConfig)
		if err != nil {
			return Repositories{}, err
		}
		mongoUserRepository, err := userrepo.NewMongoUserRepository(mongoUserRepoConfig, logger)
		if err != nil {
			return Repositories{}, err
		}
		userRepository = &mongoUserRepository
	default:
		return Repositories{}, fmt.Errorf("invalid user repository type: %s", config.UserRepository.Type)
	}
	repositories.UserRepository = userRepository

	return repositories, nil
}
