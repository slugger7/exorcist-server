package videoRepository

import (
	"database/sql"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/environment"
	"github.com/slugger7/exorcist/internal/repository/util"
)

type VideoStatement struct {
	postgres.Statement
	db *sql.DB
}

type IVideoRepository interface {
	GetVideoWithoutChecksumStatement() *VideoStatement
	UpdateVideoChecksum(video model.Video) *VideoStatement
	UpdateVideoExistsStatement(video model.Video) *VideoStatement
	GetVideosInLibraryPath(libraryPathId uuid.UUID) *VideoStatement
	InsertVideosStatement(videos []model.Video) *VideoStatement
}

type VideoRepository struct {
	db  *sql.DB
	Env *environment.EnvironmentVariables
}

var videoRepoInstance *VideoRepository

func New(db *sql.DB, env *environment.EnvironmentVariables) IVideoRepository {
	if videoRepoInstance != nil {
		return videoRepoInstance
	}
	videoRepoInstance = &VideoRepository{
		db:  db,
		Env: env,
	}
	return videoRepoInstance
}

func (ds *VideoRepository) GetVideoWithoutChecksumStatement() *VideoStatement {
	selectStatement := table.Video.SELECT(table.Video.ID, table.Video.Checksum, table.Video.RelativePath, table.LibraryPath.Path).
		FROM(table.Video.INNER_JOIN(table.LibraryPath, table.LibraryPath.ID.EQ(table.Video.LibraryPathID))).
		WHERE(table.Video.Checksum.IS_NULL())

	util.DebugCheck(ds.Env, selectStatement)

	return &VideoStatement{selectStatement, ds.db}
}

func (ds *VideoRepository) UpdateVideoChecksum(video model.Video) *VideoStatement {
	statement := table.Video.UPDATE().
		SET(table.Video.Checksum.SET(postgres.String(*video.Checksum))).
		MODEL(video).
		WHERE(table.Video.ID.EQ(postgres.UUID(video.ID)))

	util.DebugCheck(ds.Env, statement)

	return &VideoStatement{statement, ds.db}
}

func (ds *VideoRepository) UpdateVideoExistsStatement(video model.Video) *VideoStatement {
	statement := table.Video.UPDATE().
		SET(table.Video.Exists.SET(postgres.Bool(video.Exists))).
		MODEL(video).
		WHERE(table.Video.ID.EQ(postgres.UUID(video.ID)))

	util.DebugCheck(ds.Env, statement)

	return &VideoStatement{statement, ds.db}
}

func (ds *VideoRepository) GetVideosInLibraryPath(libraryPathId uuid.UUID) *VideoStatement {
	statement := table.Video.SELECT(table.Video.RelativePath, table.Video.ID).
		FROM(table.Video.Table).
		WHERE(table.Video.LibraryPathID.EQ(postgres.UUID(libraryPathId)).
			AND(table.Video.Exists.IS_TRUE()),
		)

	util.DebugCheck(ds.Env, statement)

	return &VideoStatement{statement, ds.db}
}

func (ds *VideoRepository) InsertVideosStatement(videos []model.Video) *VideoStatement {
	if len(videos) == 0 {
		return nil
	}
	statement := table.Video.INSERT(
		table.Video.LibraryPathID,
		table.Video.RelativePath,
		table.Video.Title,
		table.Video.FileName,
		table.Video.Height,
		table.Video.Width,
		table.Video.Runtime,
		table.Video.Size,
	).
		MODELS(videos)

	util.DebugCheck(ds.Env, statement)

	return &VideoStatement{statement, ds.db}
}

func (vs *VideoStatement) Query(destination interface{}) error {
	return vs.Statement.Query(vs.db, destination)
}

func (vs *VideoStatement) Exec() (sql.Result, error) {
	return vs.Statement.Exec(vs.db)
}
