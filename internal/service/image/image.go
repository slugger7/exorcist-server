package imageService

import (
	"github.com/slugger7/exorcist/internal/environment"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/repository"
)

type IImageService interface {
}

type ImageService struct {
	Env    *environment.EnvironmentVariables
	repo   repository.IRepository
	logger logger.ILogger
}

var imageServiceInstance *ImageService

func New(repo repository.IRepository, env *environment.EnvironmentVariables) IImageService {
	if imageServiceInstance == nil {
		imageServiceInstance = &ImageService{
			Env:    env,
			repo:   repo,
			logger: logger.New(env),
		}
		imageServiceInstance.logger.Info("ImageService instance created")
	}

	return imageServiceInstance
}
