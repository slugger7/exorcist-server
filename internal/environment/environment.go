package environment

import (
	"log"
	"os"
	"strconv"

	errs "github.com/slugger7/exorcist/internal/errors"
)

type ApplicationEnvironment string

var AppEnvEnum = &struct {
	Local ApplicationEnvironment
}{
	Local: "local",
}

type EnvironmentVariables struct {
	DatabaseHost     string
	DatabasePort     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string
	DebugSql         bool
	AppEnv           ApplicationEnvironment
	MediaPath        string
	Port             int
}

const (
	DATABASE_HOST     = "DATABASE_HOST"
	DATABASE_PORT     = "DATABASE_PORT"
	DATABASE_USER     = "DATABASE_USER"
	DATABASE_PASSWORD = "DATABASE_PASSWORD"
	DATABASE_NAME     = "DATABASE_NAME"
	DEBUG_SQL         = "DEBUG_SQL"
	MEDIUA_PATH       = "MEDIA_PATH"
	APP_ENV           = "APP_ENV"
	PORT              = "PORT"
)

var env *EnvironmentVariables

func GetEnvironmentVariables() *EnvironmentVariables {
	if env != nil {
		return env
	}
	RefreshEnvironmentVariables()

	return env
}

func RefreshEnvironmentVariables() {
	env = &EnvironmentVariables{
		DatabaseHost:     os.Getenv(DATABASE_HOST),
		DatabasePort:     os.Getenv(DATABASE_PORT),
		DatabaseUser:     os.Getenv(DATABASE_USER),
		DatabasePassword: os.Getenv(DATABASE_PASSWORD),
		DatabaseName:     os.Getenv(DATABASE_NAME),
		DebugSql:         getBoolValue(DEBUG_SQL, false),
		MediaPath:        os.Getenv(MEDIUA_PATH),
		AppEnv:           handleAppEnv(os.Getenv(APP_ENV)),
		Port:             getIntValue(PORT),
	}
}

func getIntValue(key string) int {
	value, err := strconv.Atoi(os.Getenv(key))
	errs.CheckError(err)

	return value
}

func getBoolValue(key string, defaultValue bool) bool {
	rawValue := os.Getenv(key)
	actualValue, err := strconv.ParseBool(rawValue)
	if err != nil {
		log.Printf("No value or invalid value found for %v setting to default '%v'\nValue was: %v", DEBUG_SQL, defaultValue, rawValue)
		return defaultValue
	}
	return actualValue
}

func handleAppEnv(value string) ApplicationEnvironment {
	switch value {
	case string(AppEnvEnum.Local):
		return AppEnvEnum.Local
	}
	log.Printf("No environment variable was set in %v defaulting to %v", APP_ENV, AppEnvEnum.Local)
	return AppEnvEnum.Local
}
