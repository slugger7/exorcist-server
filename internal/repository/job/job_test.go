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
