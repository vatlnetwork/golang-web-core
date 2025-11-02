package main

import (
	"flag"
	"golang-web-core/controllers"
	"golang-web-core/logging"
	"golang-web-core/routes"
	"golang-web-core/service"
	"golang-web-core/services/httpserver"
	"golang-web-core/terminal"
)

func main() {
	httpServerConfigPath := flag.String("http-server-config", "configs/http-server-config.json", "The path to the http server config file")
	flag.Parse()

	httpServerConfig, err := httpserver.ConfigFromJson(*httpServerConfigPath)
	if err != nil {
		panic(err)
	}

	logger := logging.NewLogger()
	logger.ServiceName = "Main"
	errorHandler, err := httpserver.NewHttpErrorHandler(&logger)
	if err != nil {
		panic(err)
	}

	applicationController, controllers, err := controllers.SetupControllers(&errorHandler)
	if err != nil {
		panic(err)
	}

	routes, err := routes.Routes(controllers, applicationController)
	if err != nil {
		panic(err)
	}

	httpServer, err := httpserver.NewHttpServer(httpServerConfig, routes, &logger)
	if err != nil {
		panic(err)
	}

	terminal := terminal.NewTerminal([]service.Service{httpServer})
	terminal.Start()
}
