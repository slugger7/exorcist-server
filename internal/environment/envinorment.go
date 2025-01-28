package environment

import (
	"log"
	"os"
	"strconv"
)

type EnvironmentVariables struct {
	DatabaseHost     string
	DatabasePort     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string
	DebugSql         bool
	Dev              bool
	MediaPath        string
}

const (
	DATABASE_HOST     = "DATABASE_HOST"
	DATABASE_PORT     = "DATABASE_PORT"
	DATABASE_USER     = "DATABASE_USER"
	DATABASE_PASSWORD = "DATABASE_PASSWORD"
	DATABASE_NAME     = "DATABASE_NAME"
	DEBUG_SQL         = "DEBUG_SQL"
	MEDIUA_PATH       = "MEDIA_PATH"
	DEV               = "DEV"
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
		Dev:              getBoolValue(DEV, false),
	}
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
