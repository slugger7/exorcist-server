package mediaService

import (
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/environment"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/models"
	"github.com/slugger7/exorcist/internal/repository"
)

type MediaService interface {
	AddPeople(id uuid.UUID, people []string) (*models.Media, error)
}

type mediaService struct {
	env    *environment.EnvironmentVariables
	repo   repository.IRepository
	logger logger.ILogger
}

// AddPeople implements MediaService.
func (m *mediaService) AddPeople(id uuid.UUID, people []string) (*models.Media, error) {
	panic("unimplemented")
}

var mediaServiceInstance *mediaService

func New(env *environment.EnvironmentVariables, repo repository.IRepository) MediaService {
	if mediaServiceInstance == nil {
		mediaServiceInstance = &mediaService{
			env:    env,
			repo:   repo,
			logger: logger.New(env),
		}

		mediaServiceInstance.logger.Info("Created media service instance")
	}

	return mediaServiceInstance
}
