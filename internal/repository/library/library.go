package libraryRepository

import (
	"github.com/go-jet/jet/v2/postgres"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/repository"
)

func CreateLibraryStatement(name string) postgres.InsertStatement {
	newLibrary := model.Library{
		Name: name,
	}

	insertStatement := table.Library.INSERT(table.Library.Name).
		MODEL(newLibrary).
		RETURNING(table.Library.ID)

	repository.DebugCheck(insertStatement)

	return insertStatement
}
