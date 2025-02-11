package libraryPathRepository

import (
	"testing"

	"github.com/slugger7/exorcist/internal/environment"
)

func Test_GetLibraryPathsSelect(t *testing.T) {
	ds := &LibraryPathRepository{
		Env: &environment.EnvironmentVariables{DebugSql: false},
	}
	statement := ds.getLibraryPathsSelect()
	sql, _ := statement.Sql()

	expectedSql := "\nSELECT library_path.id AS \"library_path.id\",\n     library_path.path AS \"library_path.path\"\nFROM public.library_path;\n"
	if sql != expectedSql {
		t.Errorf("Expected %v but got %v", expectedSql, sql)
	}
}
