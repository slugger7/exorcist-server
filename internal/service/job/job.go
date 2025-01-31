package jobService

import (
	"log"

	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/environment"
)

type IJobService interface {
	DoSomething() model.Job
}

type JobService struct {
	env environment.EnvironmentVariables
}

func New(env *environment.EnvironmentVariables) IJobService {
	jobServiceInstance := &JobService{
		env: *env,
	}

	return jobServiceInstance
}

func (js *JobService) DoSomething() model.Job {
	log.Println("Doing something")
	return model.Job{}
}
