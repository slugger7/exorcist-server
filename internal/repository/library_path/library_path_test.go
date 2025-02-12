package libraryPathRepository

import (
	"testing"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
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

func Test_Create(t *testing.T) {
	ds := &LibraryPathRepository{
		Env: &environment.EnvironmentVariables{DebugSql: false},
	}
	expectedPath := "/expected/path"
	libraryId, _ := uuid.NewRandom()
	statement := ds.create(&model.LibraryPath{Path: expectedPath, LibraryID: libraryId})
	sql, _ := statement.Sql()

	expectedSql := "\nINSERT INTO public.library_path (library_id, path)\nVALUES ($1, $2)\nRETURNING library_path.id AS \"library_path.id\",\n          library_path.path AS \"library_path.path\";\n"
	if sql != expectedSql {
		t.Errorf("Expected %v but got %v", expectedSql, sql)
	}
}
