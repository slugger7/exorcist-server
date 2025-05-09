package environment

import (
	"log"
	"os"
	"strconv"
	"strings"

	errs "github.com/slugger7/exorcist/internal/errors"
)

type ApplicationEnvironment string

var AppEnvEnum = &struct {
	Local ApplicationEnvironment
	Prod  ApplicationEnvironment
}{
	Local: "local",
	Prod:  "prod",
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
	Secret           string
	LogLevel         string
	Assets           string
	Web              *string
	JobRunner        bool
	CorsOrigins      []string
}

type OsEnv = string

const (
	DATABASE_HOST     OsEnv = "DATABASE_HOST"
	DATABASE_PORT     OsEnv = "DATABASE_PORT"
	DATABASE_USER     OsEnv = "DATABASE_USER"
	DATABASE_PASSWORD OsEnv = "DATABASE_PASSWORD"
	DATABASE_NAME     OsEnv = "DATABASE_NAME"
	DEBUG_SQL         OsEnv = "DEBUG_SQL"
	MEDIUA_PATH       OsEnv = "MEDIA_PATH"
	APP_ENV           OsEnv = "APP_ENV"
	PORT              OsEnv = "PORT"
	SECRET            OsEnv = "SECRET"
	LOG_LEVEL         OsEnv = "LOG_LEVEL"
	ASSETS            OsEnv = "ASSETS"
	WEB               OsEnv = "WEB"
	JOB_RUNNER        OsEnv = "JOB_RUNNER"
	CORS_ORIGINS      OsEnv = "CORS_ORIGINS"
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
		Secret:           os.Getenv(SECRET),
		LogLevel:         getValueOrDefault(LOG_LEVEL, "debug"),
		Assets:           os.Getenv(ASSETS),
		Web:              getValueOrNil(WEB),
		JobRunner:        getBoolValue(JOB_RUNNER, true),
		CorsOrigins:      strings.Split(os.Getenv(CORS_ORIGINS), ";"),
	}
}

func getValueOrNil(key OsEnv) *string {
	val := os.Getenv(key)
	if val == "" {
		return nil
	}
	return &val
}

func getValueOrDefault(key OsEnv, value string) string {
	val := os.Getenv(key)
	if val == "" {
		return value
	}
	return val
}

func getIntValue(key OsEnv) int {
	value, err := strconv.Atoi(os.Getenv(key))
	errs.PanicError(err)

	return value
}

func getBoolValue(key OsEnv, defaultValue bool) bool {
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
	case string(AppEnvEnum.Prod):
		return AppEnvEnum.Prod
	}
	log.Printf("No environment variable was set in %v defaulting to %v", APP_ENV, AppEnvEnum.Local)
	return AppEnvEnum.Local
}
