package util

import "fmt"

func PrintHelp() {
	fmt.Println("Please note that all of these arguments are optional. Defaults will be set for every argument if they aren't provided")
	fmt.Println("Command Line Arguments:")
	fmt.Println("   -p: (default: 3000) choose the port your application runs on")
	fmt.Println("   --cert-path: (default: none) provide the path for your ssl certificate")
	fmt.Println("   --key-path: (default: none) provide the path for your ssl certificate key")
	fmt.Println("   --disable-public-fs: if this argument is provided, the public file server will be disabled")
	fmt.Println("   --help: display this help text")
	fmt.Println("")
	fmt.Println("All environment variables are optional. Please note that environment variables are overridden by command line arguments")
	fmt.Println("Environment Variables:")
	fmt.Println("   GWC_PORT: specify a port for your application to run on")
	fmt.Println("   GWC_CERT_PATH: specify an ssl certificate file path")
	fmt.Println("   GWC_KEY_PATH: specify an ssl certificate key file path")
	fmt.Println("   GWC_ENABLE_PUBLIC_FS: [false, no] specify one of the available values to disable the public file server")
}
