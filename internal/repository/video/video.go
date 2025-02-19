package videoRepository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
)

type IVideoRepository interface {
	GetAll() ([]model.Video, error)
	GetByLibraryPathId(id uuid.UUID) ([]model.Video, error)
	GetById(id uuid.UUID) (*model.Video, error)
	UpdateExists(video *model.Video) error
	UpdateChecksum(video *model.Video) error
	Insert(models []model.Video) ([]model.Video, error)
	GetByIdWithLibraryPath(id uuid.UUID) (*VideoLibraryPathModel, error)
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
	statement := table.Video.SELECT(table.Video.AllColumns).
		FROM(table.Video)

	var vids []struct{ model.Video }
	if err := statement.Query(vr.db, &vids); err != nil {
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

func (i *VideoRepository) UpdateExists(v *model.Video) error {
	v.Modified = time.Now()
	_, err := i.updateVideoExistsStatement(*v).Exec()
	if err != nil {
		return errs.BuildError(err, "could not update video exists: %v", v.ID)
	}
	return nil
}

func (ds *VideoRepository) Insert(models []model.Video) ([]model.Video, error) {
	if len(models) == 0 {
		return nil, nil
	}

	var vids []struct{ model.Video }

	if err := ds.insertStatement(models).Query(&vids); err != nil {
		return nil, errs.BuildError(err, "could not insert video models to database")
	}

	var vidModels = []model.Video{}
	for _, v := range vids {
		vidModels = append(vidModels, v.Video)
	}

	return vidModels, nil
}

func (ds *VideoRepository) GetById(id uuid.UUID) (*model.Video, error) {
	var vids []struct{ model.Video }
	if err := ds.getByIdStatement(id).Query(&vids); err != nil {
		return nil, errs.BuildError(err, "error getting video from db for id %v", id)
	}

	var video *model.Video
	if len(vids) == 1 {
		video = &vids[len(vids)-1].Video
	}

	return video, nil
}

type VideoLibraryPathModel struct {
	model.Video
	model.LibraryPath
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

func (ds *VideoRepository) UpdateChecksum(video *model.Video) error {
	video.Modified = time.Now()
	if _, err := ds.updateChecksumStatement(*video).Exec(); err != nil {
		return errs.BuildError(err, "error updating video (%v) checksum %v", video.ID, video.Checksum)
	}

	return nil
}
