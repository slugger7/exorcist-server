package libraryPathRepository

import (
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/repository/util"
)

func (ds *LibraryPathRepository) getLibraryPathsSelect() LibraryPathStatement {
	selectQuery := table.LibraryPath.
		SELECT(table.LibraryPath.ID, table.LibraryPath.Path).
		FROM(table.LibraryPath)

	util.DebugCheck(ds.Env, selectQuery)
	return LibraryPathStatement{selectQuery, ds.db}
}

func (ds *LibraryPathRepository) create(libPath *model.LibraryPath) LibraryPathStatement {
	insertStatement := table.LibraryPath.
		INSERT(
			table.LibraryPath.LibraryID,
			table.LibraryPath.Path,
		).
		MODEL(libPath).
		RETURNING(table.LibraryPath.ID, table.LibraryPath.Path)

	util.DebugCheck(ds.Env, insertStatement)

	return LibraryPathStatement{insertStatement, ds.db}
}
