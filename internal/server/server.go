package server

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/slugger7/exorcist/internal/environment"
	"github.com/slugger7/exorcist/internal/job"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/repository"
	"github.com/slugger7/exorcist/internal/service"
)

type Server struct {
	env     *environment.EnvironmentVariables
	repo    repository.IRepository
	service service.IService
	logger  logger.ILogger
	jobCh   chan bool
	wg      *sync.WaitGroup
}

func (s *Server) withJobRunner() *Server {
	ch := job.New(s.env, s.service, s.repo, s.logger, s.wg)
	s.jobCh = ch

	ch <- true // start if any jobs exist

	return s
}

func NewServer(env *environment.EnvironmentVariables, wg *sync.WaitGroup) *http.Server {
	repo := repository.New(env)
	serv := service.New(repo, env)
	lg := logger.New(env)

	newServer := &Server{
		repo:    repo,
		env:     env,
		service: serv,
		logger:  lg,
		wg:      wg,
	}

	newServer.withJobRunner()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", env.Port),
		Handler:      newServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	server.RegisterOnShutdown(func() {
		close(newServer.jobCh)
	})

	return server
}
