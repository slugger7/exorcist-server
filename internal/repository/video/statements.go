package videoRepository

import (
	"context"
	"database/sql"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/repository/util"
)

type VideoStatement struct {
	postgres.Statement
	db  *sql.DB
	ctx context.Context
}

func (vs *VideoStatement) Query(destination interface{}) error {
	return vs.Statement.QueryContext(vs.ctx, vs.db, destination)
}

func (vs *VideoStatement) Exec() (sql.Result, error) {
	return vs.Statement.ExecContext(vs.ctx, vs.db)
}

func (ds *VideoRepository) updateChecksumStatement(video model.Video) *VideoStatement {
	statement := table.Video.UPDATE().
		SET(
			table.Video.Checksum.SET(postgres.String(*video.Checksum)),
			table.Video.Modified.SET(postgres.TimestampT(video.Modified)),
		).
		MODEL(video).
		WHERE(table.Video.ID.EQ(postgres.UUID(video.ID)))

	util.DebugCheck(ds.Env, statement)

	return &VideoStatement{statement, ds.db, ds.ctx}
}

func (ds *VideoRepository) updateVideoExistsStatement(video model.Video) *VideoStatement {
	statement := table.Video.UPDATE().
		SET(
			table.Video.Exists.SET(postgres.Bool(video.Exists)),
			table.Video.Modified.SET(postgres.TimestampT(video.Modified)),
		).
		MODEL(video).
		WHERE(table.Video.ID.EQ(postgres.UUID(video.ID)))

	util.DebugCheck(ds.Env, statement)

	return &VideoStatement{statement, ds.db, ds.ctx}
}

func (ds *VideoRepository) getByLibraryPathIdStatement(libraryPathId uuid.UUID) *VideoStatement {
	statement := table.Video.SELECT(table.Video.RelativePath, table.Video.ID).
		FROM(table.Video.Table).
		WHERE(table.Video.LibraryPathID.EQ(postgres.UUID(libraryPathId)).
			AND(table.Video.Exists.IS_TRUE()),
		)

	util.DebugCheck(ds.Env, statement)

	return &VideoStatement{statement, ds.db, ds.ctx}
}

func (ds *VideoRepository) getByIdWithLibraryPathStatement(id uuid.UUID) *VideoStatement {
	statement := table.Video.SELECT(table.Video.AllColumns, table.LibraryPath.AllColumns).
		FROM(table.Video.INNER_JOIN(table.LibraryPath, table.Video.LibraryPathID.EQ(table.LibraryPath.ID))).
		WHERE(table.Video.ID.EQ(postgres.UUID(id)).
			AND(table.Video.Deleted.IS_FALSE()).
			AND(table.Video.Exists.IS_TRUE()))

	util.DebugCheck(ds.Env, statement)

	return &VideoStatement{statement, ds.db, ds.ctx}
}
