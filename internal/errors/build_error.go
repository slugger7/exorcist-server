package errs

import (
	"errors"
	"fmt"
	"log"
	"runtime"
)

type IError interface {
	Error() string
	Message() string
	FuncName() string
	Line() int
	File() string
}

type errorDetails struct {
	file     string
	line     int
	funcName string
	message  string
}

func (e *errorDetails) Error() string {
	return fmt.Sprintf("%v@%v: %v", e.funcName, e.line, e.message)
}

func (e *errorDetails) Message() string {
	return e.message
}

func (e *errorDetails) FuncName() string {
	return e.funcName
}

func (e *errorDetails) Line() int {
	return e.line
}

func (e *errorDetails) File() string {
	return e.file
}

func BuildError(err error, format string, args ...any) error {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		log.Println("ERROR: runtime.Caller() faled")
	}
	funcName := runtime.FuncForPC(pc).Name()
	message := fmt.Sprintf(format, args...)

	e := errorDetails{
		line:     line,
		funcName: funcName,
		message:  message,
		file:     file,
	}

	return errors.Join(&e, err)
}
