package libraryRepository_test

import (
	"testing"

	libraryRepository "github.com/slugger7/exorcist/internal/repository/library"
)

func Test_CreateLibraryStatment(t *testing.T) {
	statment := libraryRepository.CreateLibraryStatement("TestName")
	sql, _ := statment.Sql()

	expectedSql := "\nINSERT INTO public.library (name)\nVALUES ($1)\nRETURNING library.id AS \"library.id\";\n"
	if sql != expectedSql {
		t.Errorf("Expected %v but got %v", expectedSql, sql)
	}
}
