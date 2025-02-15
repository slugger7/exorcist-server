package libraryService

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/repository"
)

type ILibraryService interface {
	Create(newLibrary model.Library) (*model.Library, error)
	GetAll() ([]model.Library, error)
	Action(id uuid.UUID, action string) error
}

type LibraryService struct {
	Env    *environment.EnvironmentVariables
	repo   repository.IRepository
	logger logger.ILogger
}

var libraryServiceInstance *LibraryService

func New(repo repository.IRepository, env *environment.EnvironmentVariables) ILibraryService {
	if libraryServiceInstance == nil {
		libraryServiceInstance = &LibraryService{
			Env:    env,
			repo:   repo,
			logger: logger.New(env),
		}

		libraryServiceInstance.logger.Info("LibraryService instance created")
	}
	return libraryServiceInstance
}

const ErrLibraryByName = "Could not fetch library by name %v"

func (i *LibraryService) Create(newLibrary model.Library) (*model.Library, error) {
	library, err := i.repo.Library().
		GetLibraryByName(newLibrary.Name)
	if err != nil {
		return nil, errs.BuildError(err, ErrLibraryByName, newLibrary.Name)
	}
	if library != nil {
		return nil, fmt.Errorf("library named %v already exists", newLibrary.Name)
	}

	library, err = i.repo.Library().
		CreateLibrary(newLibrary.Name)
	if err != nil {
		return nil, errs.BuildError(err, "could not create library with name %v", newLibrary.Name)
	}

	return library, nil
}

const ErrGetLibraries = "could not getting libraries in repo"

func (i *LibraryService) GetAll() ([]model.Library, error) {
	libraries, err := i.repo.Library().GetLibraries()
	if err != nil {
		return nil, errs.BuildError(err, ErrGetLibraries)
	}

	return libraries, nil
}

func (i *LibraryService) Action(id uuid.UUID, action string) error {
	panic("not implemented")
}
