package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/slugger7/exorcist/internal/environment"
	"github.com/slugger7/exorcist/internal/repository"
	"github.com/slugger7/exorcist/internal/service"
)

type Server struct {
	env     *environment.EnvironmentVariables
	repo    repository.IRepository
	service service.IService
}

func NewServer(env *environment.EnvironmentVariables) *http.Server {
	repo := repository.New(env)
	NewServer := &Server{
		repo:    repo,
		env:     env,
		service: service.New(repo, env),
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", env.Port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	return server
}
