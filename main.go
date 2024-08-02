package main

import (
	"golang-web-core/srv"
	"golang-web-core/srv/cfg"
	"golang-web-core/util"
)

func main() {
	cfg, err := cfg.FromArgs()
	if err != nil {
		util.LogFatal(err)
	}
	srv, err := srv.NewServer(cfg)
	if err != nil {
		util.LogFatal(err)
	}

	if err := srv.Start(); err != nil {
		util.LogFatal(err)
	}
}
