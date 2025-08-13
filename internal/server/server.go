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
	"github.com/slugger7/exorcist/internal/websockets"
)

type server struct {
	env       *environment.EnvironmentVariables
	repo      repository.Repository
	service   service.Service
	logger    logger.Logger
	jobCh     chan bool
	wsService websockets.Websockets
}

func (s *server) withJobRunner(ctx context.Context, wg *sync.WaitGroup, ws websockets.Websockets) *server {
	ch := job.New(s.env, s.service, s.logger, ctx, wg, ws)
	s.jobCh = ch

	ch <- true // start if any jobs exist

	return s
}

func New(env *environment.EnvironmentVariables, wg *sync.WaitGroup) *http.Server {
	lg := logger.New(env)
	shutdownCtx, cancel := context.WithCancel(context.Background())

	repo := repository.New(env, shutdownCtx)

	newServer := &server{
		repo:      repo,
		env:       env,
		logger:    lg,
		wsService: websockets.New(env),
	}

	err := newServer.repo.Job().CancelInprogress()
	if err != nil {
		lg.Errorf("clearing in progress jobs on startup: %v", err.Error())
	}

	if env.JobRunner {
		newServer.withJobRunner(shutdownCtx, wg, newServer.wsService)
	}
	newServer.service = service.New(repo, env, newServer.jobCh, shutdownCtx)

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

		newServer.logger.Debug("Cancelled and closed")

	})

	return server
}
