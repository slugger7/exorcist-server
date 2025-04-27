package imageService

import (
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/repository"
)

type IImageService interface {
	GetById(uuid.UUID) (*model.Image, error)
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

func (i *ImageService) GetById(id uuid.UUID) (*model.Image, error) {
	img, err := i.repo.Image().GetById(id)
	if err != nil {
		return nil, errs.BuildError(err, "could not get image by id: %v", id)
	}

	return img, nil
}
