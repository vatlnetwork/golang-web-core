package cfg

import (
	"golang-web-core/util"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func GetArg(key string) *string {
	for i, arg := range os.Args[1:] {
		if arg == key {
			return &os.Args[i+1]
		}
	}
	return nil
}

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
	env := os.Getenv("GWC_ENV")

	if env == "prod" || env == "production" {
		config.Env = Production
	}

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

	if GetArg("-e") != nil {
		if *GetArg("-e") == "prod" || *GetArg("-e") == "production" {
			config.Env = Production
		}
	}
	if GetArg("--help") != nil {
		util.PrintHelp()
		os.Exit(0)
	}
	if GetArg("-p") != nil {
		port := *GetArg("-p")
		portNum, err := strconv.Atoi(port)
		if err == nil {
			config.Port = portNum
		}
	}
	if GetArg("--cert-path") != nil {
		err := config.SSL.SetCertPath(*GetArg("--cert-path"))
		if err != nil {
			return config, err
		}
	}
	if GetArg("--key-path") != nil {
		err := config.SSL.SetKeyPath(*GetArg("--key-path"))
		if err != nil {
			return config, err
		}
	}
	if GetArg("--disable-public-fs") != nil {
		config.PublicFS = false
	}

	return config, nil
}
