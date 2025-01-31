package libraryPathRepository_test

import (
	"testing"

	"github.com/slugger7/exorcist/internal/environment"
	libraryPathRepository "github.com/slugger7/exorcist/internal/repository/library_path"
)

func Test_GetLibraryPathsSelect(t *testing.T) {
	ds := &libraryPathRepository.LibraryPathRepository{
		Env: &environment.EnvironmentVariables{DebugSql: false},
	}
	statement := ds.GetLibraryPathsSelect()
	sql, _ := statement.Sql()

	expectedSql := "\nSELECT library_path.id AS \"library_path.id\",\n     library_path.path AS \"library_path.path\"\nFROM public.library_path;\n"
	if sql != expectedSql {
		t.Errorf("Expected %v but got %v", expectedSql, sql)
	}
}
