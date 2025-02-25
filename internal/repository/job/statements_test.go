package jobRepository

import (
	"testing"

	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/environment"
)

var s = &JobRepository{
	Env: &environment.EnvironmentVariables{DebugSql: false},
}

func Test_CreateAllStatement(t *testing.T) {
	jobs := []model.Job{{}}

	sql, _ := s.createAllStatement(jobs).Sql()

	expected := "\nINSERT INTO public.job (job_type, status, data)\nVALUES ($1, $2, $3)\nRETURNING job.id AS \"job.id\";\n"
	if sql != expected {
		t.Errorf("Expected sql: %v\nGot sql: %v", expected, sql)
	}
}

func Test_GetNextJobStatement(t *testing.T) {
	sql, _ := s.getNextJobStatement().Sql()

	expected := "\nSELECT job.id AS \"job.id\",\n     job.status AS \"job.status\",\n     job.data AS \"job.data\",\n     job.created AS \"job.created\",\n     job.modified AS \"job.modified\",\n     job.job_type AS \"job.job_type\",\n     job.outcome AS \"job.outcome\"\nFROM public.job\nWHERE job.status = 'not_started'\nORDER BY job.created ASC\nLIMIT $1;\n"
	if sql != expected {
		t.Errorf("Expected sql: %v\nGot sql: %v", expected, sql)
	}
}
