package service

import (
	"github.com/slugger7/exorcist/internal/environment"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/repository"
	libraryService "github.com/slugger7/exorcist/internal/service/library"
	libraryPathService "github.com/slugger7/exorcist/internal/service/library_path"
	userService "github.com/slugger7/exorcist/internal/service/user"
)

type IService interface {
	UserService() userService.IUserService
	LibraryService() libraryService.ILibraryService
	LibraryPathService() libraryPathService.ILibraryPathService
}

type Service struct {
	env                *environment.EnvironmentVariables
	logger             logger.ILogger
	userService        userService.IUserService
	libraryService     libraryService.ILibraryService
	libraryPathService libraryPathService.ILibraryPathService
}

var serviceInstance *Service

func New(repo repository.IRepository, env *environment.EnvironmentVariables) IService {
	if serviceInstance == nil {
		serviceInstance = &Service{
			env:                env,
			logger:             logger.New(env),
			userService:        userService.New(repo, env),
			libraryService:     libraryService.New(repo, env),
			libraryPathService: libraryPathService.New(repo, env),
		}

		serviceInstance.logger.Info("Service instance created")
	}
	return serviceInstance
}

func (s *Service) UserService() userService.IUserService {
	s.logger.Debug("Getting UserService")
	return s.userService
}

func (s *Service) LibraryService() libraryService.ILibraryService {
	s.logger.Debug("Getting LibraryService")
	return s.libraryService
}

func (s *Service) LibraryPathService() libraryPathService.ILibraryPathService {
	s.logger.Debug("Getting LibraryPathService")
	return s.libraryPathService
}
