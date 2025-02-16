package libraryPathRepository

import (
	"testing"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/environment"
)

var ds = &LibraryPathRepository{
	Env: &environment.EnvironmentVariables{DebugSql: false},
}

func Test_GetLibraryPathsSelect(t *testing.T) {
	statement := ds.getLibraryPathsSelect()
	sql, _ := statement.Sql()

	expectedSql := "\nSELECT library_path.id AS \"library_path.id\",\n     library_path.path AS \"library_path.path\"\nFROM public.library_path;\n"
	if sql != expectedSql {
		t.Errorf("Expected %v but got %v", expectedSql, sql)
	}
}

func Test_Create(t *testing.T) {
	expectedPath := "/expected/path"
	libraryId, _ := uuid.NewRandom()
	statement := ds.create(&model.LibraryPath{Path: expectedPath, LibraryID: libraryId})
	sql, _ := statement.Sql()

	expectedSql := "\nINSERT INTO public.library_path (library_id, path)\nVALUES ($1, $2)\nRETURNING library_path.id AS \"library_path.id\",\n          library_path.path AS \"library_path.path\";\n"
	if sql != expectedSql {
		t.Errorf("Expected %v but got %v", expectedSql, sql)
	}
}

func Test_GetByLibraryId(t *testing.T) {
	id, _ := uuid.NewRandom()

	sql, _ := ds.getByLibraryIdStatement(id).Sql()

	expectedSql := "\nSELECT library_path.id AS \"library_path.id\",\n     library_path.library_id AS \"library_path.library_id\",\n     library_path.path AS \"library_path.path\",\n     library_path.created AS \"library_path.created\",\n     library_path.modified AS \"library_path.modified\"\nFROM public.library_path\nWHERE library_path.library_id = $1;\n"
	if sql != expectedSql {
		t.Errorf("Expected sql: %v\nGot sql: %v", expectedSql, sql)
	}
}

func Test_GetById(t *testing.T) {
	id, _ := uuid.NewRandom()

	sql, _ := ds.getByIdStatement(id).Sql()

	expectedSql := "\nSELECT job.id AS \"job.id\",\n     job.job_type AS \"job.job_type\",\n     job.status AS \"job.status\",\n     job.data AS \"job.data\",\n     job.created AS \"job.created\",\n     job.modified AS \"job.modified\"\nFROM public.library_path\nWHERE library_path.id = $1\nLIMIT $2;\n"
	if sql != expectedSql {
		t.Errorf("Expected sql: %v\nGot sql: %v", expectedSql, sql)
	}
}
