package personService

import (
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
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
	person, err := p.repo.Person().GetByName(name)
	if err != nil {
		return nil, errs.BuildError(err, "could not get person by name from repo")
	}

	if person == nil {
		people, err := p.repo.Person().Create([]string{name})
		if err != nil {
			return nil, errs.BuildError(err, "could not create person by name")
		}
		person = &people[0]
	}

	return person, nil
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
