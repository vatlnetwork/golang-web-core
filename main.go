package main

import (
	"inventory-app/srv"
	"inventory-app/srv/cfg"
	"inventory-app/util"
)

func main() {
	file := "default"
	_, cfgFlag := cfg.GetArg("--config")
	if cfgFlag != "" {
		file = cfgFlag
	}

	cfg, err := cfg.FromFile(file)
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
