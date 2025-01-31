package jobRepository_test

import (
	"testing"

	jobRepository "github.com/slugger7/exorcist/internal/repository/job"
)

func Test_FetchNextJob(t *testing.T) {
	actual, _ := jobRepository.FetchNextJob().Sql()

	expected := "\nSELECT job.id AS \"job.id\",\n     job.job_type AS \"job.job_type\",\n     job.status AS \"job.status\",\n     job.data AS \"job.data\",\n     job.created AS \"job.created\",\n     job.modified AS \"job.modified\"\nFROM public.job\nWHERE job.status = 'not_started'\nORDER BY job.created ASC\nLIMIT $1;\n"
	if expected != actual {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}
