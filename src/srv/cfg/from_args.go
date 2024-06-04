package cfg

import (
	"fmt"
	"golang-web-core/src/interface/database/mongo"
	"os"
	"strconv"
)

func ConfigFromArgs() (ServerConfig, error) {
	config := DefaultConfig()

	serverPort := os.Getenv("VWEB_SERVER_PORT")
	mongoDbHost := os.Getenv("VWEB_MONGO_HOST")
	mongoDbPort := os.Getenv("VWEB_MONGO_PORT")
	mongoDbUser := os.Getenv("VWEB_MONGO_USER")
	mongoDbPass := os.Getenv("VWEB_MONGO_PASS")
	mongoDbName := os.Getenv("VWEB_MONGO_DBNAME")
	certPath := os.Getenv("VWEB_CERT_PATH")
	certKeyPath := os.Getenv("VWEB_CERT_KEY_PATH")

	if serverPort != "" {
		port, err := strconv.Atoi(serverPort)
		if err != nil {
			return ServerConfig{}, err
		}
		config.Port = port
	}
	if mongoDbHost != "" {
		config.Mongo.Host = mongoDbHost
	}
	if mongoDbPort != "" {
		port, err := strconv.Atoi(mongoDbPort)
		if err != nil {
			return ServerConfig{}, err
		}
		config.Mongo.Port = port
	}
	config.Mongo.Username = mongoDbUser
	config.Mongo.Password = mongoDbPass
	if mongoDbName != "" {
		config.Mongo.DbName = mongoDbName
	}

	args := os.Args[1:]
	for i, arg := range args {
		switch arg {
		case "-p":
			port, err := strconv.Atoi(args[i+1])
			if err != nil {
				return ServerConfig{}, err
			}
			config.Port = port
		case "-e":
			env := args[i+1]
			if env == "prod" {
				config.Mongo.DbName = mongo.ProdDbName
				config.Env = Prod
			}
			if env == "m" {
				config.Mongo.DbName = mongo.TestDbName
				config.Env = MockDev
			}
		case "-ssl":
			if certPath == "" {
				return ServerConfig{}, fmt.Errorf("you must specify a certificate path in the VWEB_CERT_PATH environment variable")
			}
			if certKeyPath == "" {
				return ServerConfig{}, fmt.Errorf("you must specify a certificate path in the VWEB_CERT_KEY_PATH environment variable")
			}
			config.SSL.CertFilePath = certPath
			config.SSL.CertKeyFilePath = certKeyPath
		case "--media-directory":
			dir := args[i+1]
			if _, err := os.Stat(dir); !os.IsNotExist(err) {
				config.Media.Directory = dir
			} else {
				return ServerConfig{}, fmt.Errorf("the specified media directory, %v, does not exist", dir)
			}
		}
	}

	return config, nil
}
