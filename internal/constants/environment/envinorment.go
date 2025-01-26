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
}

const (
	DATABASE_HOST     = "DATABASE_HOST"
	DATABASE_PORT     = "DATABASE_PORT"
	DATABASE_USER     = "DATABASE_USER"
	DATABASE_PASSWORD = "DATABASE_PASSWORD"
	DATABASE_NAME     = "DATABASE_NAME"
	DEBUG_SQL         = "DEBUG_SQL"
)

func GetEnvironmentVariables() EnvironmentVariables {
	rawDebugSql := os.Getenv(DEBUG_SQL)
	debugSql, err := strconv.ParseBool(rawDebugSql)
	if err != nil {
		log.Printf("No value or invalid value found for %v setting to default 'false'\nValue was: %v", DEBUG_SQL, rawDebugSql)
	}
	envVars := EnvironmentVariables{
		DatabaseHost:     os.Getenv(DATABASE_HOST),
		DatabasePort:     os.Getenv(DATABASE_PORT),
		DatabaseUser:     os.Getenv(DATABASE_USER),
		DatabasePassword: os.Getenv(DATABASE_PASSWORD),
		DatabaseName:     os.Getenv(DATABASE_NAME),
		DebugSql:         debugSql,
	}

	return envVars
}
