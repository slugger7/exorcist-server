package libraryPathRepository

import (
	"database/sql"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/repository"
)

func GetLibraryPathsSelect() postgres.SelectStatement {
	selectQuery := table.LibraryPath.
		SELECT(table.LibraryPath.ID, table.LibraryPath.Path).
		FROM(table.LibraryPath)

	repository.DebugCheck(selectQuery)
	return selectQuery
}

func ExecuteSelect(db *sql.DB, statement postgres.SelectStatement) (data []struct{ model.LibraryPath }, err error) {
	err = statement.Query(db, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
