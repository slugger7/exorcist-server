package main

import (
	"github.com/slugger7/exorcist/internal/environment"
	"github.com/slugger7/exorcist/internal/logger"
)

func somethingNew(log logger.ILogger) {
	log.Info("Something new text")
}

func main() {
	env := environment.EnvironmentVariables{}
	logg := logger.New(&env)

	logg.Info("Who it is")
	somethingNew(logg)
}
