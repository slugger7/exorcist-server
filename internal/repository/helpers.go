package repository

import (
	"log"
	"runtime"

	"github.com/go-jet/jet/v2/postgres"
)

func (s *DatabaseService) DebugCheck(statement postgres.Statement) {
	if s.Env.DebugSql {
		pc := make([]uintptr, 10) // at least 1 entry needed
		runtime.Callers(2, pc)
		f := runtime.FuncForPC(pc[0])
		log.Printf("[%v]: %v\n", f.Name(), statement.DebugSql())
	}
}
