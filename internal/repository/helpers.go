package repository

import (
	"log"
	"runtime"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/slugger7/exorcist/internal/constants/environment"
)

func DebugCheck(statement postgres.Statement) {
	env := environment.GetEnvironmentVariables()
	if env.DebugSql {
		pc := make([]uintptr, 10) // at least 1 entry needed
		runtime.Callers(2, pc)
		f := runtime.FuncForPC(pc[0])
		log.Printf("[%v]: %v\n", f.Name(), statement.DebugSql())
	}
}
