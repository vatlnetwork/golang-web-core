package cfg

import (
	"golang-web-core/src/interface/database/mongo"
	media_os_int "golang-web-core/src/interface/media/os"
)

func DefaultConfig() ServerConfig {
	return ServerConfig{
		Port:  3001,
		Mongo: mongo.DefaultConfig(),
		Env:   IntegDev,
		Media: media_os_int.DefaultConfig(),
	}
}

func TestConfig() ServerConfig {
	cfg := DefaultConfig()
	cfg.Mongo = mongo.TestConfig()
	return cfg
}
