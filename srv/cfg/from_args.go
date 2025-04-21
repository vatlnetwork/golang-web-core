package cfg

import (
	"golang-web-core/util"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func FromArgs() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		return Config{}, err
	}

	config := Default()

	// check to see if there are environment variables
	port := os.Getenv("GWC_PORT")
	certPath := os.Getenv("GWC_CERT_PATH")
	keyPath := os.Getenv("GWC_KEY_PATH")
	enablePublicFS := os.Getenv("GWC_ENABLE_PUBLIC_FS")

	// port
	if port != "" {
		portNum, err := strconv.Atoi(port)
		if err == nil {
			config.Port = portNum
		}
	}

	// ssl
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

	// public fs
	if strings.ToLower(enablePublicFS) == "false" || strings.ToLower(enablePublicFS) == "no" {
		config.PublicFS = false
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
	}

	return config, nil
}
