package routes

import "golang-web-core/srv/cfg"

type Router struct {
	config cfg.Config
}

func NewRouter(c cfg.Config) Router {
	return Router{
		config: c,
	}
}

func (r Router) Routes() []Route {
	return []Route{}
}
