package libraryRepository

import (
	"testing"

	"github.com/slugger7/exorcist/internal/environment"
)

var lr = LibraryRepository{
	Env: &environment.EnvironmentVariables{DebugSql: false},
}

func Test_CreateLibraryStatment(t *testing.T) {
	statment := lr.createLibraryStatement("TestName")
	sql := statment.Sql()

	expectedSql := "\nINSERT INTO public.library (name)\nVALUES ($1)\nRETURNING library.id AS \"library.id\";\n"
	if sql != expectedSql {
		t.Errorf("Expected %v but got %v", expectedSql, sql)
	}
}

func Test_GetLibraryByName(t *testing.T) {
	statment := lr.getLibraryByNameStatement("TestName")
	sql := statment.Sql()

	expectedSql := "\nSELECT library.id AS \"library.id\"\nFROM public.library\nWHERE library.name = $1::text;\n"
	if sql != expectedSql {
		t.Errorf("Expected %v but got %v", expectedSql, sql)
	}
}

func Test_GetLibraries(t *testing.T) {
	statment := lr.getLibrariesStatement()
	sql := statment.Sql()

	expectedSql := "\nSELECT library.id AS \"library.id\",\n     library.name AS \"library.name\"\nFROM public.library;\n"
	if sql != expectedSql {
		t.Errorf("Expected %v but got %v", expectedSql, sql)
	}
}
