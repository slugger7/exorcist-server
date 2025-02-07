package logger

import (
	"log"
	"os"
	"runtime"

	"github.com/slugger7/exorcist/internal/environment"
)

type ILogger interface {
	Debug(message string)
	Info(message string)
	Warning(message string)
	Error(message string)
}

type logger struct {
	env           *environment.EnvironmentVariables
	debugLogger   *log.Logger
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
}

var loggerInstance *logger

func New(env *environment.EnvironmentVariables) ILogger {
	if loggerInstance == nil {
		loggerInstance = &logger{
			env:           env,
			debugLogger:   log.New(os.Stdout, "[DEBUG]", log.LUTC),
			infoLogger:    log.New(os.Stdout, "[INFO] ", log.Default().Flags()),
			warningLogger: log.New(os.Stdout, "[WARN]", 0),
			errorLogger:   log.New(os.Stdout, "[ERROR]", 0),
		}
	}
	return loggerInstance
}

func reflectFunction() (string, string, int) {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[1])
	file, line := f.FileLine(pc[1])
	return file, f.Name(), line
}

func (l *logger) Debug(message string) {
	file, funcName, line := reflectFunction()
	l.debugLogger.Printf("%v@%v(%v): %v", file, line, funcName, message)
}

func (l *logger) Info(message string) {
	l.infoLogger.Println(message)
}

func (l *logger) Warning(message string) {
	_, funcName, _ := reflectFunction()
	l.debugLogger.Printf("%v: %v", funcName, message)
}

func (l *logger) Error(message string) {
	file, funcName, line := reflectFunction()
	l.errorLogger.Printf("%v@%v(%v): %v", file, line, funcName, message)
}
