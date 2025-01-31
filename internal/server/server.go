package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/slugger7/exorcist/internal/db"
	"github.com/slugger7/exorcist/internal/environment"
)

type Server struct {
	port int
	db   db.Service
	env  *environment.EnvironmentVariables
}

func NewServer(env *environment.EnvironmentVariables) *http.Server {
	NewServer := &Server{
		port: env.Port,
		db:   db.New(env),
		env:  env,
	}

	if err := NewServer.db.RunMigrations(); err != nil {
		log.Printf("Colud not run migrations because: %v", err)
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	return server
}
