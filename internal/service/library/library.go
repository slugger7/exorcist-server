package libraryService

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/dto"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/models"
	"github.com/slugger7/exorcist/internal/repository"
)

type LibraryService interface {
	Create(newLibrary *model.Library) (*model.Library, error)
	GetAll() ([]model.Library, error)
	GetMedia(id, userId uuid.UUID, search dto.MediaSearchDTO) (*dto.PageDTO[models.MediaOverviewModel], error)
}

type libraryService struct {
	env    *environment.EnvironmentVariables
	repo   repository.IRepository
	logger logger.ILogger
}

// GetMedia implements LibraryService.
func (i *libraryService) GetMedia(id, userId uuid.UUID, search dto.MediaSearchDTO) (*dto.PageDTO[models.MediaOverviewModel], error) {
	library, err := i.repo.Library().GetById(id)
	if err != nil {
		return nil, errs.BuildError(err, "could not get library by id from repo: %v", id)
	}

	if library == nil {
		return nil, fmt.Errorf("no library found with id: %v", id.String())
	}

	media, err := i.repo.Library().GetMedia(id, userId, search)
	if err != nil {
		return nil, errs.BuildError(err, "could not get media for library (%v) from repo", id.String())
	}

	return media, nil
}

var libraryServiceInstance *libraryService

func New(repo repository.IRepository, env *environment.EnvironmentVariables) LibraryService {
	if libraryServiceInstance == nil {
		libraryServiceInstance = &libraryService{
			env:    env,
			repo:   repo,
			logger: logger.New(env),
		}

		libraryServiceInstance.logger.Info("LibraryService instance created")
	}
	return libraryServiceInstance
}

const (
	ErrLibraryByName string = "could not fetch library by name %v"
	ErrLibraryExists string = "library named %v already exists"
)

func (i *libraryService) Create(newLibrary *model.Library) (*model.Library, error) {
	library, err := i.repo.Library().
		GetByName(newLibrary.Name)
	if err != nil {
		return nil, errs.BuildError(err, ErrLibraryByName, newLibrary.Name)
	}
	if library != nil {
		return nil, fmt.Errorf(ErrLibraryExists, newLibrary.Name)
	}

	library, err = i.repo.Library().
		Create(newLibrary.Name)
	if err != nil {
		return nil, errs.BuildError(err, "could not create library with name %v", newLibrary.Name)
	}

	return library, nil
}

const ErrGetLibraries = "could not getting libraries in repo"

func (i *libraryService) GetAll() ([]model.Library, error) {
	libraries, err := i.repo.Library().GetAll()
	if err != nil {
		return nil, errs.BuildError(err, ErrGetLibraries)
	}

	return libraries, nil
}
