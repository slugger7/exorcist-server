package personRepository

import (
	"context"
	"database/sql"

	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/environment"
	"github.com/slugger7/exorcist/internal/logger"
)

type IPersonRepository interface {
	GetByName(name string) (*model.Person, error)
	Create(name string) (*model.Person, error)
}

type personRepository struct {
	env    *environment.EnvironmentVariables
	db     *sql.DB
	logger logger.ILogger
	ctx    context.Context
}

// Create implements IPersonRepository.
func (p *personRepository) Create(name string) (*model.Person, error) {
	panic("unimplemented")
}

// GetByName implements IPersonRepository.
func (p *personRepository) GetByName(name string) (*model.Person, error) {
	panic("unimplemented")
}

var personRepositoryInstance *personRepository

func New(env *environment.EnvironmentVariables, db *sql.DB, context context.Context) IPersonRepository {
	if personRepositoryInstance == nil {
		personRepositoryInstance = &personRepository{
			env:    env,
			db:     db,
			logger: logger.New(env),
			ctx:    context,
		}

		personRepositoryInstance.logger.Info("PersonRepository instance created")
	}

	return personRepositoryInstance
}
