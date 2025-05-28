package service

import (
	"github.com/slugger7/exorcist/internal/environment"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/repository"
	jobService "github.com/slugger7/exorcist/internal/service/job"
	libraryService "github.com/slugger7/exorcist/internal/service/library"
	libraryPathService "github.com/slugger7/exorcist/internal/service/library_path"
	userService "github.com/slugger7/exorcist/internal/service/user"
)

type IService interface {
	User() userService.IUserService
	Library() libraryService.ILibraryService
	LibraryPath() libraryPathService.ILibraryPathService
	Job() jobService.IJobService
}

type Service struct {
	env         *environment.EnvironmentVariables
	logger      logger.ILogger
	user        userService.IUserService
	library     libraryService.ILibraryService
	libraryPath libraryPathService.ILibraryPathService
	job         jobService.IJobService
}

var serviceInstance *Service

func New(repo repository.IRepository, env *environment.EnvironmentVariables, jobCh chan bool) IService {
	if serviceInstance == nil {
		serviceInstance = &Service{
			env:         env,
			logger:      logger.New(env),
			user:        userService.New(repo, env),
			library:     libraryService.New(repo, env),
			libraryPath: libraryPathService.New(repo, env),
			job:         jobService.New(repo, env, jobCh),
		}

		serviceInstance.logger.Info("Service instance created")
	}
	return serviceInstance
}

func (s *Service) User() userService.IUserService {
	s.logger.Debug("Getting UserService")
	return s.user
}

func (s *Service) Library() libraryService.ILibraryService {
	s.logger.Debug("Getting LibraryService")
	return s.library
}

func (s *Service) LibraryPath() libraryPathService.ILibraryPathService {
	s.logger.Debug("Getting LibraryPathService")
	return s.libraryPath
}

func (s *Service) Job() jobService.IJobService {
	s.logger.Debug("Getting jobService")
	return s.job
}
