package videoService

import (
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/repository"
)

type IVideoService interface {
	GetAll() ([]model.Video, error)
	GetById(id uuid.UUID) (*model.Video, error)
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
