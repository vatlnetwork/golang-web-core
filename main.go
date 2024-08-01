package main

import (
	"golang-web-core/srv"
	"golang-web-core/srv/cfg"
	"log"
)

func main() {
	cfg, err := cfg.FromArgs()
	if err != nil {
		log.Fatal(err)
	}
	srv, err := srv.NewServer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
