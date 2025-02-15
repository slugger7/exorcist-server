package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/slugger7/exorcist/internal/environment"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/repository"
	"github.com/slugger7/exorcist/internal/service"
)

type Server struct {
	env     *environment.EnvironmentVariables
	repo    repository.IRepository
	service service.IService
	logger  logger.ILogger
}

func NewServer(env *environment.EnvironmentVariables) *http.Server {
	repo := repository.New(env)
	newServer := &Server{
		repo:    repo,
		env:     env,
		service: service.New(repo, env),
		logger:  logger.New(env),
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", env.Port),
		Handler:      newServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	return server
}
