package main

import (
	"golang-web-core/src/srv"
	"golang-web-core/src/srv/cfg"
)

func main() {
	config, err := cfg.ConfigFromArgs()
	if err != nil {
		panic(err)
	}
	srv, err := srv.NewServer(config)
	if err != nil {
		panic(err)
	}
	err = srv.Run()
	if err != nil {
		panic(err)
	}
}
