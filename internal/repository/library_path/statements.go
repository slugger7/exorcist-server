package libraryPathRepository

import (
	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/repository/util"
)

func (ds *LibraryPathRepository) getLibraryPathsSelect() LibraryPathStatement {
	selectQuery := table.LibraryPath.
		SELECT(table.LibraryPath.ID, table.LibraryPath.Path).
		FROM(table.LibraryPath)

	util.DebugCheck(ds.Env, selectQuery)
	return LibraryPathStatement{selectQuery, ds.db, ds.ctx}
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

	return LibraryPathStatement{insertStatement, ds.db, ds.ctx}
}

func (lps *LibraryPathRepository) getByLibraryIdStatement(libraryId uuid.UUID) LibraryPathStatement {
	statement := table.LibraryPath.SELECT(table.LibraryPath.AllColumns).
		FROM(table.LibraryPath).
		WHERE(table.LibraryPath.LibraryID.EQ(postgres.UUID(libraryId)))

	util.DebugCheck(lps.Env, statement)

	return LibraryPathStatement{statement, lps.db, lps.ctx}
}

func (lps *LibraryPathRepository) getByIdStatement(id uuid.UUID) LibraryPathStatement {
	statement := table.LibraryPath.SELECT(table.LibraryPath.AllColumns).
		FROM(table.LibraryPath).
		WHERE(table.LibraryPath.ID.EQ(postgres.UUID(id))).
		LIMIT(1)

	util.DebugCheck(lps.Env, statement)

	return LibraryPathStatement{statement, lps.db, lps.ctx}
}
