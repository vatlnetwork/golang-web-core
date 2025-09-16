package controllers

import (
	"errors"
	"golang-web-core/services/httpserver"
)

func SetupControllers(errorHandler *httpserver.HttpErrorHandler) (httpserver.Controller, []httpserver.Controller, error) {
	if errorHandler == nil {
		return nil, nil, errors.New("error handler is required")
	}

	applicationController, err := NewApplicationController(errorHandler)
	if err != nil {
		return nil, nil, err
	}

	return applicationController, []httpserver.Controller{}, nil
}
