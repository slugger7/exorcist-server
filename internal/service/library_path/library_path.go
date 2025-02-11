package libraryPathService

import (
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/environment"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/repository"
)

type ILibraryPathService interface {
	Create(*model.LibraryPath) (*model.LibraryPath, error)
}

type LibraryPathService struct {
	Env    *environment.EnvironmentVariables
	repo   repository.IRepository
	logger logger.ILogger
}

var libraryPathServiceInstance *LibraryPathService

func New(repo repository.IRepository, env *environment.EnvironmentVariables) ILibraryPathService {
	if libraryPathServiceInstance == nil {
		libraryPathServiceInstance = &LibraryPathService{
			Env:    env,
			repo:   repo,
			logger: logger.New(env),
		}

		libraryPathServiceInstance.logger.Info("LibraryPathService instance created")
	}

	return libraryPathServiceInstance
}

func (lps *LibraryPathService) Create(libPathModel *model.LibraryPath) (*model.LibraryPath, error) {
	panic("not implemented")
}
