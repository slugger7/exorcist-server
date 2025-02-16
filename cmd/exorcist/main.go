package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/server"
)

func gracefulShutdown(apiServer *http.Server, done chan bool, wg *sync.WaitGroup) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Waiting for job runner to finish")
	wg.Wait()

	log.Println("Server exiting")

	done <- true
}

func main() {
	err := godotenv.Load()
	errs.PanicError(err)
	env := environment.GetEnvironmentVariables()

	var wg sync.WaitGroup
	server := server.NewServer(env, &wg)

	done := make(chan bool, 1)

	go gracefulShutdown(server, done, &wg)

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	<-done
	log.Println("Graceful shutdown complete.")
}
