package service

import (
	"log"

	"github.com/slugger7/exorcist/internal/environment"
	"github.com/slugger7/exorcist/internal/repository"
	userService "github.com/slugger7/exorcist/internal/service/user"
)

type IService interface {
	UserService() userService.IUserService
}

type Service struct {
	env         *environment.EnvironmentVariables
	userService userService.IUserService
}

var serviceInstance *Service

func New(repo repository.IRepository, env *environment.EnvironmentVariables) IService {
	if serviceInstance == nil {
		serviceInstance = &Service{
			env:         env,
			userService: userService.New(repo, env),
		}

		log.Println("Service instance created")
	}
	return serviceInstance
}

func (s *Service) UserService() userService.IUserService {
	log.Println("Getting UserService")
	return s.userService
}
