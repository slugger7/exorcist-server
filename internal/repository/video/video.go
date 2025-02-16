package videoRepository

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
)

type IVideoRepository interface {
	GetAll() ([]model.Video, error)
	GetByLibraryPathId(id uuid.UUID) ([]model.Video, error)
	UpdateVideoExists(video model.Video) error
	Insert(models []model.Video) error
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

func (vr *VideoRepository) GetAll() ([]model.Video, error) {
	panic("not implemented")
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

func (i *VideoRepository) UpdateVideoExists(v model.Video) error {
	_, err := i.updateVideoExistsStatement(v).Exec()
	if err != nil {
		return errs.BuildError(err, "could not update video exists: %v", v.ID)
	}
	return nil
}

func (ds *VideoRepository) Insert(models []model.Video) error {
	if len(models) == 0 {
		return nil
	}

	if _, err := ds.insertStatement(models).Exec(); err != nil {
		return errs.BuildError(err, "could not insert video models to database")
	}

	return nil
}
