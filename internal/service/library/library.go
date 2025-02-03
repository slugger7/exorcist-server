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
	var libraries []struct {
		model.Library
	}
	if err := i.repo.LibraryRepo().
		GetLibraryByName(newLibrary.Name).
		Query(&libraries); err != nil {
		log.Printf("Could not fetch library by name %v", newLibrary.Name)
		return nil, err
	}
	if len(libraries) > 0 {
		return nil, errors.New(fmt.Sprintf("library named %v already exists", newLibrary.Name))
	}

	if err := i.repo.LibraryRepo().
		CreateLibraryStatement(newLibrary.Name).
		Query(&libraries); err != nil {
		log.Printf("could not create library with name %v", newLibrary.Name)
		return nil, err
	}

	return &libraries[len(libraries)-1].Library, nil
}
