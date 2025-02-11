package libraryPathService

import (
	"fmt"

	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/repository"
)

const (
	LibraryPathWasNilErr = "library path model was nil"
	LibraryNilErr        = "library was nil for id: %v"
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
	if libPathModel == nil {
		return nil, fmt.Errorf(LibraryPathWasNilErr)
	}

	library, err := lps.repo.LibraryRepo().GetLibraryById(libPathModel.LibraryID)
	if err != nil {
		return nil, errs.BuildError(err, "could not get library by id")
	}

	if library == nil {
		return nil, fmt.Errorf(LibraryNilErr, libPathModel.LibraryID)
	}

	libraryPath, err := lps.repo.LibraryPathRepo().Create(libPathModel)
	if err != nil {
		return nil, errs.BuildError(err, "could not create new library path")
	}

	return libraryPath, nil
}
