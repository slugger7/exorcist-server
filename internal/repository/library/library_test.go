package libraryRepository_test

import (
	"testing"

	"github.com/slugger7/exorcist/internal/environment"
	libraryRepository "github.com/slugger7/exorcist/internal/repository/library"
)

func Test_CreateLibraryStatment(t *testing.T) {
	ds := &libraryRepository.LibraryRepository{
		Env: &environment.EnvironmentVariables{DebugSql: false},
	}
	statment := ds.CreateLibraryStatement("TestName")
	sql, _ := statment.Sql()

	expectedSql := "\nINSERT INTO public.library (name)\nVALUES ($1)\nRETURNING library.id AS \"library.id\";\n"
	if sql != expectedSql {
		t.Errorf("Expected %v but got %v", expectedSql, sql)
	}
}
