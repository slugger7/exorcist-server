package media

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/repository/util"
)

type IMediaRepository interface {
	Create([]model.Media) ([]model.Media, error)
	UpdateExists(model.Media) error
}

type MediaRepository struct {
	db     *sql.DB
	Env    *environment.EnvironmentVariables
	logger logger.ILogger
	ctx    context.Context
}

var mediaRepositoryInstance *MediaRepository

func New(db *sql.DB, env *environment.EnvironmentVariables, context context.Context) IMediaRepository {
	if mediaRepositoryInstance != nil {
		return mediaRepositoryInstance
	}

	mediaRepositoryInstance = &MediaRepository{
		db:     db,
		Env:    env,
		ctx:    context,
		logger: logger.New(env),
	}

	return mediaRepositoryInstance
}

func (r *MediaRepository) Create(ms []model.Media) ([]model.Media, error) {
	if len(ms) == 0 {
		return nil, nil
	}
	media := table.Media

	statement := media.INSERT(
		media.LibraryPathID,
		media.Path,
		media.FileName,
		media.Title,
		media.Size,
	).
		MODELS(ms).
		RETURNING(media.AllColumns)

	util.DebugCheck(r.Env, statement)

	models := []model.Media{}

	if err := statement.Query(r.db, &models); err != nil {
		return nil, errs.BuildError(err, "could not insert media models")
	}

	return models, nil
}

func (r *MediaRepository) UpdateExists(m model.Media) error {
	m.Modified = time.Now()

	statement := table.Video.UPDATE().
		SET(
			table.Media.Exists.SET(postgres.Bool(m.Exists)),
			table.Media.Modified.SET(postgres.TimestampT(m.Modified)),
		).
		MODEL(m).
		WHERE(table.Video.ID.EQ(postgres.UUID(m.ID)))

	util.DebugCheck(r.Env, statement)

	if _, err := statement.Exec(r.db); err != nil {
		return errs.BuildError(err, "could not update media exists: %v", m.ID)
	}

	return nil
}

func (r *MediaRepository) UpdateChecksum(m model.Media) error {
	m.Modified = time.Now()
	statement := table.Video.UPDATE().
		SET(
			table.Media.Checksum.SET(postgres.String(*m.Checksum)),
			table.Media.Modified.SET(postgres.TimestampT(m.Modified)),
		).
		MODEL(m).
		WHERE(table.Video.ID.EQ(postgres.UUID(m.ID)))

	util.DebugCheck(r.Env, statement)

	if _, err := statement.Exec(r.db); err != nil {
		return errs.BuildError(err, "could not update checksum for video: %v", m.ID)
	}

	return nil
}
