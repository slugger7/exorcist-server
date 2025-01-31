package service

import (
	"github.com/slugger7/exorcist/internal/environment"
	jobService "github.com/slugger7/exorcist/internal/service/job"
)

type Service struct {
	env        *environment.EnvironmentVariables
	JobService jobService.IJobService
}

func New(env *environment.EnvironmentVariables) *Service {
	serviceInstance := &Service{
		env:        env,
		JobService: jobService.New(env),
	}
	return serviceInstance
}
