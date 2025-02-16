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
