package videoRepository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/repository/util"
)

type VideoLibraryPathModel struct {
	model.Video
	model.LibraryPath
}

type IVideoRepository interface {
	GetAll() ([]model.Video, error)
	GetByLibraryPathId(id uuid.UUID) ([]model.Video, error)
	Insert(models []model.Video) ([]model.Video, error)
	GetByIdWithLibraryPath(id uuid.UUID) (*VideoLibraryPathModel, error)
}

type VideoRepository struct {
	db     *sql.DB
	Env    *environment.EnvironmentVariables
	logger logger.ILogger
	ctx    context.Context
}

var videoRepoInstance *VideoRepository

func New(db *sql.DB, env *environment.EnvironmentVariables, context context.Context) IVideoRepository {
	if videoRepoInstance != nil {
		return videoRepoInstance
	}
	videoRepoInstance = &VideoRepository{
		db:     db,
		Env:    env,
		ctx:    context,
		logger: logger.New(env),
	}
	return videoRepoInstance
}

func (vr *VideoRepository) GetAll() ([]model.Video, error) {
	statement := table.Video.SELECT(table.Video.AllColumns).
		FROM(table.Video)

	var vids []struct{ model.Video }
	if err := statement.QueryContext(vr.ctx, vr.db, &vids); err != nil {
		return nil, errs.BuildError(err, "could not get all videos")
	}

	if vids == nil {
		return nil, nil
	}

	var models []model.Video
	for _, v := range vids {
		models = append(models, v.Video)
	}

	return models, nil
}

func (ds *VideoRepository) GetByLibraryPathId(id uuid.UUID) ([]model.Video, error) {
	var vids []struct{ model.Video }
	if err := ds.getByLibraryPathIdStatement(id).Query(&vids); err != nil {
		return nil, errs.BuildError(err, "could not get videos by library path id: %v", id)
	}

	vidModels := []model.Video{}
	for _, v := range vids {
		vidModels = append(vidModels, v.Video)
	}

	return vidModels, nil
}

func (r *VideoRepository) Insert(models []model.Video) ([]model.Video, error) {
	if len(models) == 0 {
		return nil, nil
	}

	statement := table.Video.INSERT(
		table.Video.MediaID,
		table.Video.Height,
		table.Video.Width,
		table.Video.Runtime,
	).
		MODELS(models).
		RETURNING(table.Video.AllColumns)

	util.DebugCheck(r.Env, statement)

	var vids []model.Video

	if err := statement.Query(r.db, &vids); err != nil {
		return nil, errs.BuildError(err, "could not insert video models to database")
	}

	return vids, nil
}

func (ds *VideoRepository) GetByIdWithLibraryPath(id uuid.UUID) (*VideoLibraryPathModel, error) {
	var results []VideoLibraryPathModel
	if err := ds.getByIdWithLibraryPathStatement(id).Query(&results); err != nil {
		return nil, errs.BuildError(err, "error getting video by id (%v) with library path", id.String())
	}

	var result VideoLibraryPathModel
	if results != nil {
		result = results[len(results)-1]
	}

	return &result, nil
}
