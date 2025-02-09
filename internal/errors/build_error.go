package errs

import (
	"fmt"
	"log"
	"runtime"
)

func BuildError(err error, format string, args ...any) error {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		log.Println("ERROR: runtime.Caller() faled")
	}
	funcName := runtime.FuncForPC(pc).Name()
	message := fmt.Sprintf(format, args...)

	return fmt.Errorf("%v: %v\n%w", funcName, message, err)
}
