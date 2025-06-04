package personService

import (
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/environment"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/repository"
)

type IPersonService interface {
	Upsert(name string) (*model.Person, error)
}

type personService struct {
	env    *environment.EnvironmentVariables
	repo   repository.IRepository
	logger logger.ILogger
}

// Upsert implements IPersonService.
func (p *personService) Upsert(name string) (*model.Person, error) {
	panic("unimplemented")
}

var personServiceInstance *personService

func New(repo repository.IRepository, env *environment.EnvironmentVariables) IPersonService {
	if personServiceInstance == nil {
		personServiceInstance = &personService{
			env:    env,
			repo:   repo,
			logger: logger.New(env),
		}

		personServiceInstance.logger.Info("PersonService instance created")
	}

	return personServiceInstance
}
