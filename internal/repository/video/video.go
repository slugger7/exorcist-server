package videoRepository

import (
	"context"
	"database/sql"

	"github.com/go-jet/jet/v2/postgres"
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

type MediaVideoModel struct {
	model.Video
	model.Media
}

type IVideoRepository interface {
	GetAll() ([]model.Video, error)
	Insert(models []model.Video) ([]model.Video, error)
	GetByIdWithMedia(id uuid.UUID) (*MediaVideoModel, error)
	GetByMediaId(id uuid.UUID) (*MediaVideoModel, error)
}

type VideoRepository struct {
	db     *sql.DB
	Env    *environment.EnvironmentVariables
	logger logger.ILogger
	ctx    context.Context
}

// GetByMediaId implements IVideoRepository.
func (vr *VideoRepository) GetByMediaId(id uuid.UUID) (*MediaVideoModel, error) {
	video := table.Video
	media := table.Media

	statement := video.SELECT(video.AllColumns, media.AllColumns).
		FROM(video.INNER_JOIN(media, video.MediaID.EQ(media.ID))).
		LIMIT(1)

	util.DebugCheck(vr.Env, statement)

	var result MediaVideoModel

	if err := statement.QueryContext(vr.ctx, vr.db, &result); err != nil {
		return nil, errs.BuildError(err, "could not get video by media id: %v", id)
	}

	return &result, nil
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

func (r *VideoRepository) GetByIdWithMedia(id uuid.UUID) (*MediaVideoModel, error) {
	video := table.Video
	media := table.Media
	statement := video.SELECT(video.AllColumns, media.AllColumns).
		FROM(video.INNER_JOIN(video, video.MediaID.EQ(media.ID))).
		WHERE(video.ID.EQ(postgres.UUID(id))).
		LIMIT(1)

	util.DebugCheck(r.Env, statement)

	var result MediaVideoModel
	if err := statement.QueryContext(r.ctx, r.db, &result); err != nil {
		return nil, errs.BuildError(err, "could not find video by id with media: %v", id)
	}

	return &result, nil
}
