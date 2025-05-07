package server

import (
	"context"
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
}

func (s *Server) withJobRunner(ctx context.Context, wg *sync.WaitGroup) *Server {
	ch := job.New(s.env, s.service, s.repo, s.logger, ctx, wg)
	s.jobCh = ch

	ch <- true // start if any jobs exist

	return s
}

func NewServer(env *environment.EnvironmentVariables, wg *sync.WaitGroup) *http.Server {
	repo := repository.New(env)
	lg := logger.New(env)
	shutdownCtx, cancel := context.WithCancel(context.Background())

	newServer := &Server{
		repo:   repo,
		env:    env,
		logger: lg,
	}

	if env.JobRunner {
		newServer.withJobRunner(shutdownCtx, wg)
	}
	newServer.service = service.New(repo, env, newServer.jobCh)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", env.Port),
		Handler:      newServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	server.RegisterOnShutdown(func() {
		newServer.logger.Info("Shutting down server. Stopping job runner.")
		cancel()
		close(newServer.jobCh)
	})

	return server
}
