package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/slugger7/exorcist/internal/environment"
	"github.com/slugger7/exorcist/internal/job"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/models"
	"github.com/slugger7/exorcist/internal/repository"
	"github.com/slugger7/exorcist/internal/service"
)

type server struct {
	env            *environment.EnvironmentVariables
	repo           repository.IRepository
	service        service.IService
	logger         logger.ILogger
	jobCh          chan bool
	websockets     models.WebSocketMap
	websocketMutex sync.Mutex
}

func (s *server) withJobRunner(ctx context.Context, wg *sync.WaitGroup, wss models.WebSocketMap) *server {
	ch := job.New(s.env, s.service, s.logger, ctx, wg, wss)
	s.jobCh = ch

	ch <- true // start if any jobs exist

	return s
}

func New(env *environment.EnvironmentVariables, wg *sync.WaitGroup) *http.Server {
	lg := logger.New(env)
	shutdownCtx, cancel := context.WithCancel(context.Background())

	repo := repository.New(env, shutdownCtx)

	newServer := &server{
		repo:       repo,
		env:        env,
		logger:     lg,
		websockets: make(models.WebSocketMap),
	}

	err := newServer.repo.Job().CancelInprogress()
	if err != nil {
		lg.Errorf("clearing in progress jobs on startup: %v", err.Error())
	}

	if env.JobRunner {
		newServer.withJobRunner(shutdownCtx, wg, newServer.websockets)
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

		newServer.websocketMutex.Lock()
		defer newServer.websocketMutex.Unlock()

		newServer.logger.Debug("Closing websockets")
		for _, i := range newServer.websockets {
			for _, s := range i {
				s.Mu.Lock()
				defer s.Mu.Unlock()
				s.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				s.Conn.Close()
			}
		}

		newServer.logger.Debug("Websockets closed")
	})

	return server
}
