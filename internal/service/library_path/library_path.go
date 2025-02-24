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
	Create(m *model.LibraryPath) (*model.LibraryPath, error)
	GetAll() ([]model.LibraryPath, error)
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

const ErrGetLibraryById = "could not get library by id: %v"
const ErrCreateLibraryPath = "could not create new library path"

func (lps *LibraryPathService) Create(libPathModel *model.LibraryPath) (*model.LibraryPath, error) {
	if libPathModel == nil {
		return nil, fmt.Errorf(LibraryPathWasNilErr)
	}

	library, err := lps.repo.Library().GetLibraryById(libPathModel.LibraryID)
	if err != nil {
		return nil, errs.BuildError(err, ErrGetLibraryById, libPathModel.LibraryID)
	}

	if library == nil {
		return nil, fmt.Errorf(LibraryNilErr, libPathModel.LibraryID)
	}

	libraryPath, err := lps.repo.LibraryPath().Create(libPathModel.Path, libPathModel.LibraryID)
	if err != nil {
		return nil, errs.BuildError(err, ErrCreateLibraryPath)
	}

	return libraryPath, nil
}

const ErrGetAllLibraryPaths = "could not get all library paths"

func (lps *LibraryPathService) GetAll() ([]model.LibraryPath, error) {
	libPaths, err := lps.repo.LibraryPath().GetAll()
	if err != nil {
		return nil, errs.BuildError(err, ErrGetAllLibraryPaths)
	}

	return libPaths, nil
}
