package cfg

import (
	"os"
	"strconv"
	"strings"
)

func GetEnvironment() Config {
	config := Development()

	// check to see if there are environment variables
	env := os.Getenv("GWC_ENV")
	if env == "prod" || env == "production" {
		config = Production()
	}

	args := os.Args

	// args override the env variables
	for i, arg := range args {
		if arg == "-e" {
			env := args[i+1]
			if env == "prod" {
				config = Production()
			}
			if env == "dev" {
				config = Development()
			}
		}
	}

	return config
}

func FromArgs() (Config, error) {
	config := GetEnvironment()

	// check to see if there are environment variables
	port := os.Getenv("GWC_PORT")
	certPath := os.Getenv("GWC_CERT_PATH")
	keyPath := os.Getenv("GWC_KEY_PATH")
	enablePublicFS := os.Getenv("GWC_ENABLE_PUBLIC_FS")

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

	args := os.Args

	// args override the env variables
	for i, arg := range args {
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
	}

	return config, nil
}
