package libraryPathRepository

import (
	"database/sql"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/repository"
)

func GetLibraryPathsSelect() postgres.SelectStatement {
	selectQuery := table.LibraryPath.
		SELECT(table.LibraryPath.ID, table.LibraryPath.Path).
		FROM(table.LibraryPath)

	repository.DebugCheckSelect(selectQuery)
	return selectQuery
}

func CreateLibraryPath(libraryId uuid.UUID, path string) postgres.InsertStatement {
	newLibPath := model.LibraryPath{
		LibraryID: libraryId,
		Path:      path,
	}

	insertStatement := table.LibraryPath.
		INSERT(
			table.LibraryPath.LibraryID,
			table.LibraryPath.Path,
		).
		MODEL(newLibPath).
		RETURNING(table.LibraryPath.ID, table.LibraryPath.Path)

	repository.DebugCheckInsert(insertStatement)

	return insertStatement
}

func ExecuteSelect(db *sql.DB, statement postgres.SelectStatement) (data []struct{ model.LibraryPath }, err error) {
	err = statement.Query(db, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func ExecuteInsert(db *sql.DB, statement postgres.InsertStatement) (data []struct{ model.LibraryPath }, err error) {
	err = statement.Query(db, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
