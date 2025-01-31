package util

import (
	"log"
	"runtime"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/slugger7/exorcist/internal/environment"
)

func DebugCheck(env *environment.EnvironmentVariables, statement postgres.Statement) {
	if env.DebugSql {
		pc := make([]uintptr, 10) // at least 1 entry needed
		runtime.Callers(2, pc)
		f := runtime.FuncForPC(pc[0])
		log.Printf("[%v]: %v\n", f.Name(), statement.DebugSql())
	}
}
