package main

import (
	"golang-web-core/controllers"
	"golang-web-core/logging"
	"golang-web-core/routes"
	"golang-web-core/service"
	"golang-web-core/services/httpserver"
	"golang-web-core/terminal"
)

func main() {
	config, err := httpserver.ConfigFromJson("configs/http-server-config.json")
	if err != nil {
		panic(err)
	}

	logger := logging.NewLogger()
	applicationController, err := controllers.NewApplicationController(&logger)
	if err != nil {
		panic(err)
	}

	routes, err := routes.Routes(nil, applicationController)
	if err != nil {
		panic(err)
	}

	httpServer, err := httpserver.NewHttpServer(config, routes, &logger)
	if err != nil {
		panic(err)
	}

	terminal := terminal.NewTerminal([]service.Service{httpServer})
	terminal.Start()
}
