package media

import (
	"context"
	"database/sql"

	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/logger"
)

type IMediaRepository interface {
	Create([]model.Media) ([]model.Media, error)
}

type MediaRepository struct {
	db     *sql.DB
	Env    *environment.EnvironmentVariables
	logger logger.ILogger
	ctx    context.Context
}

var mediaRepositoryInstance *MediaRepository

func New(db *sql.DB, env *environment.EnvironmentVariables, context context.Context) IVideoRepository {
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

	models := []model.Media{}

	if err := statement.Query(r.db, &models); err != nil {
		return nil, errs.BuildError(err, "could not insert media models")
	}

	return models, nil
}
