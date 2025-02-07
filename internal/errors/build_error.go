package errs

import (
	"errors"
	"fmt"
	"log"
	"runtime"
)

func BuildError(err error, format string, args ...any) error {
	pc, _, lineNo, ok := runtime.Caller(1)
	if !ok {
		log.Println("ERROR: runtime.Caller() faled")
	}
	funcName := runtime.FuncForPC(pc).Name()
	message := fmt.Sprintf(format, args...)

	return errors.Join(fmt.Errorf("%v@%v: %v", funcName, lineNo, message), err)
}
