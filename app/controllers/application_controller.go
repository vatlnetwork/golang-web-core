package controllers

import (
	"fmt"
	"golang-web-core/srv/cfg"
	"log"
	"net/http"
)

// you shouldn't be touching this file except for the BeforeAction and setupControllers

type ApplicationController struct {
	cfg.Config
	Controllers map[string]Controller
}

func NewApplicationController(config cfg.Config) (ApplicationController, error) {
	cont := ApplicationController{
		Config:      config,
		Controllers: map[string]Controller{},
	}

	err := cont.setupControllers()

	return cont, err
}

func (c ApplicationController) Name() string {
	return "ApplicationController"
}

func (c ApplicationController) BeforeAction(handler http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		// any checks you want to do on every single request that goes into the server can go here
		handler(rw, req)
	}
}

func (c ApplicationController) setupControllers() error {
	controllers := []Controller{
		// this is where you initialize your controllers. if you do not initialize your controllers here, they will not be usable
		NewTestController(c.Config),
	}

	// everything below here should be left untouched

	for _, cont := range controllers {
		_, ok := c.Controllers[cont.Name()]
		if ok {
			return fmt.Errorf("error: a controller with the name %v was registered twice", cont.Name())
		}
		c.Controllers[cont.Name()] = cont
	}

	return nil
}

func (c ApplicationController) GetController(name string) Controller {
	controller, ok := c.Controllers[name]
	if !ok {
		log.Fatalf("attempted to access a controller that does not exist! %v", name)
	}

	return controller
}
