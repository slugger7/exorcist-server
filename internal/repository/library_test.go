package repository_test

import (
	"testing"

	"github.com/slugger7/exorcist/internal/environment"
	"github.com/slugger7/exorcist/internal/repository"
)

func Test_CreateLibraryStatment(t *testing.T) {
	ds := &repository.DatabaseService{
		Env: &environment.EnvironmentVariables{DebugSql: false},
	}
	statment := ds.CreateLibraryStatement("TestName")
	sql, _ := statment.Sql()

	expectedSql := "\nINSERT INTO public.library (name)\nVALUES ($1)\nRETURNING library.id AS \"library.id\";\n"
	if sql != expectedSql {
		t.Errorf("Expected %v but got %v", expectedSql, sql)
	}
}
