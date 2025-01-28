package main

import (
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/slugger7/exorcist/internal/db"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/job"
)

func main() {
	err := godotenv.Load()
	errs.CheckError(err)

	env := environment.GetEnvironmentVariables()

	db := db.NewDatabase(env)

	job.GenerateChecksums(db)
}
