package libraryRepository

import (
	"context"
	"database/sql"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/repository/util"
)

type LibraryStatement struct {
	postgres.Statement
	db  *sql.DB
	ctx context.Context
}

func (ls *LibraryRepository) createLibraryStatement(name string) *LibraryStatement {
	newLibrary := model.Library{
		Name: name,
	}

	insertStatement := table.Library.INSERT(table.Library.Name).
		MODEL(newLibrary).
		RETURNING(table.Library.ID)

	util.DebugCheck(ls.Env, insertStatement)

	return &LibraryStatement{insertStatement, ls.db, ls.ctx}
}

func (i *LibraryRepository) getLibraryByNameStatement(name string) *LibraryStatement {
	statement := table.Library.SELECT(table.Library.ID).
		FROM(table.Library).
		WHERE(table.Library.Name.EQ(postgres.String(name)))

	util.DebugCheck(i.Env, statement)
	return &LibraryStatement{statement, i.db, i.ctx}
}

func (ls *LibraryRepository) getLibrariesStatement() *LibraryStatement {
	statement := table.Library.SELECT(table.Library.AllColumns).
		FROM(table.Library)

	util.DebugCheck(ls.Env, statement)

	return &LibraryStatement{statement, ls.db, ls.ctx}
}

func (ls *LibraryRepository) getById(id uuid.UUID) *LibraryStatement {
	statement := table.Library.SELECT(table.Library.ID, table.Library.Name).
		FROM(table.Library).
		WHERE(table.Library.ID.EQ(postgres.UUID(id)))

	util.DebugCheck(ls.Env, statement)

	return &LibraryStatement{statement, ls.db, ls.ctx}
}
