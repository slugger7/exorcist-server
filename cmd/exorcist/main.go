package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/slugger7/exorcist/internal/db"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/job"
)

func main() {
	err := godotenv.Load()
	errs.CheckError(err)
	env := environment.GetEnvironmentVariables()

	database := db.NewDatabase(env)
	defer database.Close()

	err = db.RunMigrations(database, env)
	if err != nil {
		log.Printf("Error occured when running migrations: %v", err.Error())
	}

	job.ScanPath(database)
}
