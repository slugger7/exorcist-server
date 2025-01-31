package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
)

func NewDatabase(env *environment.EnvironmentVariables) *sql.DB {
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		env.DatabaseHost,
		env.DatabasePort,
		env.DatabaseUser,
		env.DatabasePassword,
		env.DatabaseName)
	if env.AppEnv == environment.AppEnvEnum.Local {
		log.Printf("connection_string: %v", psqlconn)
	}
	db, err := sql.Open("postgres", psqlconn)
	errs.CheckError(err)

	return db
}

func RunMigrations(db *sql.DB, env *environment.EnvironmentVariables) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return err
	}

	log.Println("Running migrations")
	err = m.Up()
	if err != nil {
		return err
	}
	return nil
}
