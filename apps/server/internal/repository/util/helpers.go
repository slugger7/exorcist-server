package util

import (
	"log"
	"runtime"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/slugger7/exorcist/internal/environment"
	"github.com/slugger7/exorcist/internal/logger"
)

func DebugCheck(env *environment.EnvironmentVariables, statement postgres.Statement) {
	logg := logger.New(env)
	if env.DebugSql {
		pc, _, lineNo, ok := runtime.Caller(1)
		if !ok {
			log.Println("DebugCheck: runtime.Caller() failed")
		}
		funcName := runtime.FuncForPC(pc).Name()

		logg.Debugf("\n%v@%v: %v", funcName, lineNo, statement.DebugSql())
	}
}
