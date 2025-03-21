package cfg

import (
	"golang-web-core/srv/database_adapters/imdb"
	"golang-web-core/srv/database_adapters/mongo"
	"golang-web-core/util"
	"os"
	"strconv"
	"strings"
)

type Environment string

const (
	Dev  Environment = "dev"
	Prod Environment = "prod"
)

func GetEnvironment() Environment {
	env := os.Getenv("GWC_ENV")
	if env == "prod" || env == "production" {
		return Prod
	}

	args := os.Args
	for i, arg := range args {
		if arg == "-e" {
			env := args[i+1]
			if env == "prod" {
				return Prod
			}
		}
	}

	return Dev
}

func GetEnvironmentConfig() Config {
	env := GetEnvironment()

	if env == Prod {
		return Production()
	}

	return Development()
}

func FromArgs() (Config, error) {
	config := GetEnvironmentConfig()

	// check to see if there are environment variables
	port := os.Getenv("GWC_PORT")
	certPath := os.Getenv("GWC_CERT_PATH")
	keyPath := os.Getenv("GWC_KEY_PATH")
	enablePublicFS := os.Getenv("GWC_ENABLE_PUBLIC_FS")
	databaseAdapter := os.Getenv("GWC_DB_ADAPTER")
	databaseHostname := os.Getenv("GWC_DB_HOSTNAME")
	databaseName := os.Getenv("GWC_DB_NAME")
	databaseUsername := os.Getenv("GWC_DB_USERNAME")
	databasePassword := os.Getenv("GWC_DB_PASSWORD")

	if port != "" {
		portNum, err := strconv.Atoi(port)
		if err == nil {
			config.Port = portNum
		}
	}
	if certPath != "" {
		err := config.SSL.SetCertPath(certPath)
		if err != nil {
			return config, err
		}
	}
	if keyPath != "" {
		err := config.SSL.SetKeyPath(keyPath)
		if err != nil {
			return config, err
		}
	}
	if strings.ToLower(enablePublicFS) == "false" || strings.ToLower(enablePublicFS) == "no" {
		config.PublicFS = false
	}
	switch databaseAdapter {
	case "imdb":
		config.Database.Adapter = imdb.NewImdbAdapter()
		config.Database.Connection = config.Database.Adapter.Connection()
	case "mongo":
		config.Database.Adapter = mongo.NewMongoAdapter()
		config.Database.Connection = config.Database.Adapter.Connection()
	case "none":
		config.Database.Adapter = nil
	}
	if databaseHostname != "" {
		config.Database.Connection.Hostname = databaseHostname
	}
	if databaseName != "" {
		config.Database.Connection.Database = databaseName
	}
	if databaseUsername != "" {
		config.Database.Connection.Username = databaseUsername
	}
	if databasePassword != "" {
		config.Database.Connection.Password = databasePassword
	}

	args := os.Args

	// args override the env variables
	for i, arg := range args {
		if arg == "--help" {
			util.PrintHelp()
			os.Exit(0)
		}
		if arg == "-p" {
			port := args[i+1]
			portNum, err := strconv.Atoi(port)
			if err == nil {
				config.Port = portNum
			}
		}
		if arg == "--cert-path" {
			err := config.SSL.SetCertPath(args[i+1])
			if err != nil {
				return config, err
			}
		}
		if arg == "--key-path" {
			err := config.SSL.SetKeyPath(args[i+1])
			if err != nil {
				return config, err
			}
		}
		if arg == "--disable-public-fs" {
			config.PublicFS = false
		}
		if arg == "--db-adapter" {
			switch args[i+1] {
			case "imdb":
				config.Database.Adapter = imdb.NewImdbAdapter()
				config.Database.Connection = config.Database.Adapter.Connection()
			case "mongo":
				config.Database.Adapter = mongo.NewMongoAdapter()
				config.Database.Connection = config.Database.Adapter.Connection()
			}
		}
		if arg == "--db-host" {
			config.Database.Connection.Hostname = args[i+1]
		}
		if arg == "--db-name" {
			config.Database.Connection.Database = args[i+1]
		}
		if arg == "--db-user" {
			config.Database.Connection.Username = args[i+1]
		}
		if arg == "--db-pass" {
			config.Database.Connection.Password = args[i+1]
		}
		if arg == "--no-db" {
			config.Database.Adapter = nil
		}
	}

	return config, nil
}
