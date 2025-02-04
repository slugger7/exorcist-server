package service

import (
	"log"

	"github.com/slugger7/exorcist/internal/environment"
	"github.com/slugger7/exorcist/internal/repository"
	libraryService "github.com/slugger7/exorcist/internal/service/library"
	userService "github.com/slugger7/exorcist/internal/service/user"
)

type IService interface {
	UserService() userService.IUserService
	LibraryService() libraryService.ILibraryService
}

type Service struct {
	env            *environment.EnvironmentVariables
	userService    userService.IUserService
	libraryService libraryService.ILibraryService
}

var serviceInstance *Service

func New(repo repository.IRepository, env *environment.EnvironmentVariables) IService {
	if serviceInstance == nil {
		serviceInstance = &Service{
			env:            env,
			userService:    userService.New(repo, env),
			libraryService: libraryService.New(repo, env),
		}

		log.Println("Service instance created")
	}
	return serviceInstance
}

func (s *Service) UserService() userService.IUserService {
	log.Println("Getting UserService")
	return s.userService
}

func (s *Service) LibraryService() libraryService.ILibraryService {
	log.Println("Getting LibraryService")
	return s.libraryService
}
