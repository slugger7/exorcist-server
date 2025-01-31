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
	repo    repository.IRepository
	env     *environment.EnvironmentVariables
	Service service.Service
}

func NewServer(env *environment.EnvironmentVariables) *http.Server {
	NewServer := &Server{
		repo:    repository.New(env),
		env:     env,
		Service: *service.New(env),
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
