package videoService

import (
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/environment"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/repository"
)

type IVideoService interface {
	GetAll() ([]model.Video, error)
}

type VideoService struct {
	Env    *environment.EnvironmentVariables
	repo   repository.IRepository
	logger logger.ILogger
}

var videoServiceInstance *VideoService

func New(repo repository.IRepository, env *environment.EnvironmentVariables) *VideoService {
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

func (vs *VideoService) GetAll() ([]model.Video, error) {
	panic("not implemented")
}
