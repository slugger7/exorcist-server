package libraryService

import (
	"errors"
	"fmt"
	"log"

	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/environment"
	"github.com/slugger7/exorcist/internal/repository"
)

type ILibraryService interface {
	CreateLibrary(newLibrary model.Library) (*model.Library, error)
}

type LibraryService struct {
	Env  *environment.EnvironmentVariables
	repo repository.IRepository
}

var libraryServiceInstance *LibraryService

func New(repo repository.IRepository, env *environment.EnvironmentVariables) *LibraryService {
	if libraryServiceInstance == nil {
		libraryServiceInstance = &LibraryService{
			Env:  env,
			repo: repo,
		}

		log.Println("LibraryService instance created")
	}
	return libraryServiceInstance
}

func (i LibraryService) CreateLibrary(newLibrary model.Library) (*model.Library, error) {
	library, err := i.repo.LibraryRepo().
		GetLibraryByName(newLibrary.Name)
	if err != nil {
		log.Printf("Could not fetch library by name %v", newLibrary.Name)
		return nil, err
	}
	if library != nil {
		return nil, fmt.Errorf("library named %v already exists", newLibrary.Name)
	}

	library, err = i.repo.LibraryRepo().
		CreateLibrary(newLibrary.Name)
	if err != nil {
		log.Printf("could not create library with name %v", newLibrary.Name)
		return nil, err
	}

	return library, nil
}

func (i LibraryService) GetLibraries() ([]model.Library, error) {
	libraries, err := i.repo.LibraryRepo().GetLibraries()
	if err != nil {
		return nil, errors.Join(errors.New("error getting libraries in repo"), err)
	}

	return libraries, nil
}
