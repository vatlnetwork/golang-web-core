package cfg

import (
	"fmt"
	"os"
	"strconv"
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

	if port != "" {
		portNum, err := strconv.Atoi(port)
		if err == nil {
			config.Port = portNum
		}
	}
	if certPath != "" {
		_, err := os.ReadFile(certPath)
		if err != nil {
			if os.IsNotExist(err) {
				return config, fmt.Errorf("the cert path you specified (%v) does not exist", certPath)
			}
			return config, err
		}
	}
	if keyPath != "" {
		_, err := os.ReadFile(keyPath)
		if err != nil {
			if os.IsNotExist(err) {
				return config, fmt.Errorf("the cert path you specified (%v) does not exist", keyPath)
			}
			return config, err
		}
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
	}

	return config, nil
}
