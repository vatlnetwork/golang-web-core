package controllers

import (
	"errors"
	"golang-web-core/domain"
	"golang-web-core/repositories"
	"golang-web-core/services/httpserver"
)

func SetupControllers(repositories repositories.Repositories, errorHandler *httpserver.HttpErrorHandler) (httpserver.Controller, []httpserver.Controller, error) {
	if errorHandler == nil {
		return nil, nil, errors.New("error handler is required")
	}

	sessionManager := domain.NewSessionManager(repositories.SessionRepository, repositories.UserRepository)

	applicationController, err := NewApplicationController(errorHandler, sessionManager)
	if err != nil {
		return nil, nil, err
	}

	controllers := []httpserver.Controller{}

	authController, err := NewAuthController(sessionManager, errorHandler)
	if err != nil {
		return nil, nil, err
	}
	controllers = append(controllers, authController)

	return applicationController, controllers, nil
}
