package main

import (
	"github.com/slugger7/exorcist/internal/environment"
	"github.com/slugger7/exorcist/internal/logger"
)

func functionWrappedLogs(logg logger.ILogger) {
	logg.Debug("func std debug")
	logg.Debugf("func fmt %v", "debug")

	logg.Info("func std info")
	logg.Infof("func fmt %v", "info")

	logg.Warning("func std warning")
	logg.Warningf("func fmt %v", "warning")

	logg.Error("func std error")
	logg.Errorf("func fmt %v", "error")
}

func main() {
	env := environment.EnvironmentVariables{LogLevel: "none"}
	logg := logger.New(&env)

	logg.Debug("std debug")
	logg.Debugf("fmt %v", "debug")

	logg.Info("std info")
	logg.Infof("fmt %v", "info")

	logg.Warning("std warning")
	logg.Warningf("fmt %v", "warning")

	logg.Error("std error")
	logg.Errorf("fmt %v", "error")

	functionWrappedLogs(logg)
}
