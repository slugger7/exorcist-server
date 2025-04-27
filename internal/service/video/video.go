package videoService

import (
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/models"
	"github.com/slugger7/exorcist/internal/repository"
	videoRepository "github.com/slugger7/exorcist/internal/repository/video"
)

type IVideoService interface {
	GetAll() ([]model.Video, error)
	GetOverview() ([]models.VideoOverviewDTO, error)
	GetById(id uuid.UUID) (*model.Video, error)
	GetByIdWithLibraryPath(id uuid.UUID) (*videoRepository.VideoLibraryPathModel, error)
}

type VideoService struct {
	Env    *environment.EnvironmentVariables
	repo   repository.IRepository
	logger logger.ILogger
}

var videoServiceInstance *VideoService

func New(repo repository.IRepository, env *environment.EnvironmentVariables) IVideoService {
	if videoServiceInstance == nil {
		videoServiceInstance = &VideoService{
			Env:    env,
			repo:   repo,
			logger: logger.New(env),
		}

		videoServiceInstance.logger.Info("VideoService instance created")
	}
	return videoServiceInstance
}

func (vs *VideoService) GetOverview() ([]models.VideoOverviewDTO, error) {
	vids, err := vs.repo.Video().GetOverview()
	if err != nil {
		return nil, errs.BuildError(err, "could not get videos for overview")
	}

	videos := make([]models.VideoOverviewDTO, len(vids))
	for i, v := range vids {
		videos[i] = models.VideoOverviewDTO{
			Id:            v.Video.ID,
			Title:         v.Video.Title,
			Path:          v.LibraryPath.Path + v.Video.RelativePath,
			ThumbnailPath: v.Image.Path,
		}
	}

	return videos, nil
}

const ErrGetAllVideos = "could not get all videos"

func (vs *VideoService) GetAll() ([]model.Video, error) {
	vids, err := vs.repo.Video().GetAll()
	if err != nil {
		return nil, errs.BuildError(err, ErrGetAllVideos)
	}

	return vids, nil
}

const ErrVideoById = "error getting video by id %v"

func (vs *VideoService) GetById(id uuid.UUID) (*model.Video, error) {
	video, err := vs.repo.Video().GetById(id)
	if err != nil {
		return nil, errs.BuildError(err, ErrVideoById, id)
	}

	return video, nil
}

func (vs *VideoService) GetByIdWithLibraryPath(id uuid.UUID) (*videoRepository.VideoLibraryPathModel, error) {
	return vs.repo.Video().GetByIdWithLibraryPath(id)
}
