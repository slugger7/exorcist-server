package personRepository

import (
	"context"
	"database/sql"
	"strings"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/repository/util"
)

var person = table.Person

type IPersonRepository interface {
	GetByName(name string) (*model.Person, error)
	Create(names []string) ([]model.Person, error)
}

type personRepository struct {
	env    *environment.EnvironmentVariables
	db     *sql.DB
	logger logger.ILogger
	ctx    context.Context
}

// Create implements IPersonRepository.
func (p *personRepository) Create(names []string) ([]model.Person, error) {
	if len(names) == 0 {
		return nil, nil
	}

	peopleModels := make([]model.Person, len(names))
	for i, n := range names {
		peopleModels[i] = model.Person{Name: n}
	}

	statement := person.INSERT(person.Name).
		MODELS(peopleModels).
		RETURNING(person.AllColumns)

	util.DebugCheck(p.env, statement)

	var createdModels []model.Person
	if err := statement.Query(p.db, &createdModels); err != nil {
		return nil, errs.BuildError(err, "could not insert new people models")
	}

	return createdModels, nil
}

// GetByName implements IPersonRepository.
func (p *personRepository) GetByName(name string) (*model.Person, error) {
	statement := person.SELECT(person.AllColumns).
		FROM(person).
		WHERE(postgres.LOWER(person.Name).EQ(postgres.String(strings.ToLower(name)))).
		LIMIT(1)

	util.DebugCheck(p.env, statement)

	var personModel []model.Person

	if err := statement.QueryContext(p.ctx, p.db, &personModel); err != nil {
		return nil, errs.BuildError(err, "could not query people by name")
	}

	if len(personModel) == 0 {
		return nil, nil
	}

	return &personModel[0], nil
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
