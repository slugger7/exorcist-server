package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/slugger7/exorcist/internal/constants/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
)

func NewDatabase(env environment.EnvironmentVariables) *sql.DB {
	if env.Dev {
		log.Printf("host=%s port=%s user=%s password=%s database=%s",
			env.DatabaseHost,
			env.DatabasePort,
			env.DatabaseUser,
			env.DatabasePassword,
			env.DatabaseName)
	}
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		env.DatabaseHost,
		env.DatabasePort,
		env.DatabaseUser,
		env.DatabasePassword,
		env.DatabaseName)
	db, err := sql.Open("postgres", psqlconn)
	errs.CheckError(err)

	return db
}
