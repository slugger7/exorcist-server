package libraryRepository_test

import (
	"testing"

	"github.com/slugger7/exorcist/internal/environment"
	libraryRepository "github.com/slugger7/exorcist/internal/repository/library"
)

var lr = libraryRepository.LibraryRepository{
	Env: &environment.EnvironmentVariables{DebugSql: false},
}

func Test_CreateLibraryStatment(t *testing.T) {
	statment := lr.CreateLibraryStatement("TestName")
	sql, _ := statment.Sql()

	expectedSql := "\nINSERT INTO public.library (name)\nVALUES ($1)\nRETURNING library.id AS \"library.id\";\n"
	if sql != expectedSql {
		t.Errorf("Expected %v but got %v", expectedSql, sql)
	}
}

func Test_GetLibraryByName(t *testing.T) {
	statment := lr.GetLibraryByName("TestName")
	sql, _ := statment.Sql()

	expectedSql := "\nSELECT library.id AS \"library.id\"\nFROM public.library\nWHERE library.name = $1::text;\n"
	if sql != expectedSql {
		t.Errorf("Expected %v but got %v", expectedSql, sql)
	}
}
