package library_test

import (
	"testing"

	repo "github.com/slugger7/exorcist/internal/repository/library"
)

func Test_CreateLibraryStatment(t *testing.T) {
	statment := repo.CreateLibraryStatement("TestName")
	sql := statment.DebugSql()

	expectedSql := "\nINSERT INTO public.library (name)\nVALUES ('TestName')\nRETURNING library.id AS \"library.id\";\n"
	if sql != expectedSql {
		t.Errorf("Expected %v but got %v", expectedSql, sql)
	}
}
