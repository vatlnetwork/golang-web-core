package util

import "fmt"

func PrintHelp() {
	fmt.Println("Please note that all of these arguments are optional. Defaults will be set for every argument if they aren't provided")
	fmt.Println("Command Line Arguments:")
	fmt.Println("   -e: [dev, prod] (default: dev) choose the environment your application runs in")
	fmt.Println("   -p: (default: 3000) choose the port your application runs on")
	fmt.Println("   --cert-path: (default: none) provide the path for your ssl certificate")
	fmt.Println("   --key-path: (default: none) provide the path for your ssl certificate key")
	fmt.Println("   --disable-public-fs: if this argument is provided, the public file server will be disabled")
	fmt.Println("   --db-adapter: [imdb, mongo] (default: imdb) choose the database your application will use (if it uses one)")
	fmt.Println("   --db-host: (default: none) provide the host and port for your database server. should be in the format host:port")
	fmt.Println("   --db-name: (default: none) provide the name of the database your application will use")
	fmt.Println("   --db-user: (default: none) provide a username for the database server your application will use")
	fmt.Println("   --db-pass: (default: none) provide a password for the database server your application will use")
	fmt.Println("   --no-db: add this argument if you do not wish to enable any databases in your application")
	fmt.Println("   --help: display this help text")
	fmt.Println("")
	fmt.Println("All environment variables are optional. Please note that environment variables are overridden by command line arguments")
	fmt.Println("Environment Variables:")
	fmt.Println("   GWC_ENV: [prod, production] specify whether to use production mode")
	fmt.Println("   GWC_PORT: specify a port for your application to run on")
	fmt.Println("   GWC_CERT_PATH: specify an ssl certificate file path")
	fmt.Println("   GWC_KEY_PATH: specify an ssl certificate key file path")
	fmt.Println("   GWC_ENABLE_PUBLIC_FS: [false, no] specify one of the available values to disable the public file server")
	fmt.Println("   GWC_DB_ADAPTER: [imdb, mongo, none] specify a database adapter to use or none to disable the database")
	fmt.Println("   GWC_DB_HOSTNAME: specify the host and port for your database server. should be in the format host:port")
	fmt.Println("   GWC_DB_NAME: specify the name of the database your application will use")
	fmt.Println("   GWC_DB_USERNAME: specify a username to connect to your database server")
	fmt.Println("   GWC_DB_PASSWORD: specify a password to connect to your database server")
}
