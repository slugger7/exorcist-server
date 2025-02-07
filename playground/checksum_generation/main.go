package main

import (
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/job"
	"github.com/slugger7/exorcist/internal/repository"
)

func main() {
	err := godotenv.Load()
	errs.PanicError(err)

	env := environment.GetEnvironmentVariables()

	db := repository.New(env)

	job.GenerateChecksums(db)
}
