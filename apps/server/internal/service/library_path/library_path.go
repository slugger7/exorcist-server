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

type LibraryPathService interface {
	Create(m *model.LibraryPath) (*model.LibraryPath, error)
	GetAll() ([]model.LibraryPath, error)
}

type libraryPathService struct {
	env    *environment.EnvironmentVariables
	repo   repository.Repository
	logger logger.Logger
}

var libraryPathServiceInstance *libraryPathService

func New(repo repository.Repository, env *environment.EnvironmentVariables) LibraryPathService {
	if libraryPathServiceInstance == nil {
		libraryPathServiceInstance = &libraryPathService{
			env:    env,
			repo:   repo,
			logger: logger.New(env),
		}

		libraryPathServiceInstance.logger.Info("LibraryPathService instance created")
	}

	return libraryPathServiceInstance
}

const ErrGetLibraryById = "could not get library by id: %v"
const ErrCreateLibraryPath = "could not create new library path"

func (lps *libraryPathService) Create(libPathModel *model.LibraryPath) (*model.LibraryPath, error) {
	if libPathModel == nil {
		return nil, fmt.Errorf(LibraryPathWasNilErr)
	}

	libPathsExist, err := lps.repo.LibraryPath().GetContainingPath(libPathModel.Path)
	if err != nil {
		return nil, errs.BuildError(err, "could not get paths containing path")
	}

	if libPathsExist != nil {
		if len(libPathsExist) > 0 {
			return nil, fmt.Errorf("found paths that contain path or contains the path")
		}
	}

	library, err := lps.repo.Library().GetById(libPathModel.LibraryID)
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

func (lps *libraryPathService) GetAll() ([]model.LibraryPath, error) {
	libPaths, err := lps.repo.LibraryPath().GetAll()
	if err != nil {
		return nil, errs.BuildError(err, ErrGetAllLibraryPaths)
	}

	return libPaths, nil
}
