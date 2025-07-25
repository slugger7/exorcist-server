package personRepository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/dto"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/models"
	"github.com/slugger7/exorcist/internal/repository/helpers"
	"github.com/slugger7/exorcist/internal/repository/util"
)

var person = table.Person
var mediaPerson = table.MediaPerson

type PersonRepository interface {
	GetById(id uuid.UUID) (*model.Person, error)
	GetByName(name string) (*model.Person, error)
	Create(names []string) ([]model.Person, error)
	AddToMedia(mediaPeople []model.MediaPerson) ([]model.MediaPerson, error)
	RemoveFromMedia(mediaPerson model.MediaPerson) error
	GetAll(search dto.PersonSearchDTO) ([]model.Person, error)
	GetMedia(id uuid.UUID, search dto.MediaSearchDTO) (*dto.PageDTO[models.MediaOverviewModel], error)
}

type personRepository struct {
	env    *environment.EnvironmentVariables
	db     *sql.DB
	logger logger.ILogger
	ctx    context.Context
}

// GetMedia implements PersonRepository.
func (r *personRepository) GetMedia(id uuid.UUID, search dto.MediaSearchDTO) (*dto.PageDTO[models.MediaOverviewModel], error) {
	search.People = []string{}
	relationFunc := func(relationTable postgres.ReadableTable) postgres.ReadableTable {
		media := table.Media
		mediaPerson := table.MediaPerson
		return relationTable.INNER_JOIN(
			mediaPerson,
			media.ID.EQ(mediaPerson.MediaID),
		)
	}

	whereFn := func(whr postgres.BoolExpression) postgres.BoolExpression {
		return whr.AND(mediaPerson.PersonID.EQ(postgres.UUID(id)))
	}

	selectStatement, countStatement := helpers.MediaOverviewStatement(search, relationFunc, whereFn)

	util.DebugCheck(r.env, countStatement)
	util.DebugCheck(r.env, selectStatement)

	var total struct {
		Total int
	}
	if err := countStatement.QueryContext(r.ctx, r.db, &total); err != nil {
		return nil, errs.BuildError(err, "could not query media for total")
	}

	var mediaResult []models.MediaOverviewModel
	if err := selectStatement.QueryContext(r.ctx, r.db, &mediaResult); err != nil {
		return nil, errs.BuildError(err, "could not query media")
	}

	return &dto.PageDTO[models.MediaOverviewModel]{
		Data:  mediaResult,
		Limit: search.Limit,
		Skip:  search.Skip,
		Total: total.Total,
	}, nil
}

// GetById implements PersonRepository.
func (p *personRepository) GetById(id uuid.UUID) (*model.Person, error) {
	statement := person.SELECT(person.AllColumns).
		WHERE(person.ID.EQ(postgres.UUID(id)))

	var peopleModels []model.Person
	if err := statement.QueryContext(p.ctx, p.db, &peopleModels); err != nil {
		return nil, errs.BuildError(err, "could not fetch person by id from db %v", id)
	}

	if len(peopleModels) == 0 {
		return nil, nil
	}

	return &peopleModels[0], nil
}

// GetAll implements PersonRepository.
func (p *personRepository) GetAll(search dto.PersonSearchDTO) ([]model.Person, error) {
	statement := person.SELECT(person.AllColumns).
		FROM(
			person.LEFT_JOIN(table.MediaPerson, person.ID.EQ(table.MediaPerson.PersonID)),
		).
		GROUP_BY(person.ID).
		ORDER_BY(search.ToOrderByClause()...)

	if search.Search != "" {
		caseInsensitive := strings.ToLower(search.Search)
		statement = statement.WHERE(postgres.LOWER(person.Name).LIKE(postgres.String(fmt.Sprintf("%%%v%%", caseInsensitive))))
	}

	var people []model.Person
	if err := statement.QueryContext(p.ctx, p.db, &people); err != nil {
		return nil, errs.BuildError(err, "could not fetch people from database")
	}

	return people, nil
}

// RemoveFromMedia implements IPersonRepository.
func (p *personRepository) RemoveFromMedia(mp model.MediaPerson) error {
	statement := mediaPerson.DELETE().
		WHERE(mediaPerson.MediaID.EQ(postgres.UUID(mp.MediaID)).
			AND(mediaPerson.PersonID.EQ(postgres.UUID(mp.PersonID))))

	util.DebugCheck(p.env, statement)

	if _, err := statement.ExecContext(p.ctx, p.db); err != nil {
		return errs.BuildError(err, "could not delete media person with media id %v and person id %v", mp.MediaID.String(), mp.PersonID.String())
	}

	return nil
}

// AddToMedia implements IPersonRepository.
func (p *personRepository) AddToMedia(mediaPeople []model.MediaPerson) ([]model.MediaPerson, error) {
	if len(mediaPeople) == 0 {
		return nil, nil
	}

	statement := mediaPerson.INSERT(mediaPerson.MediaID, mediaPerson.PersonID).
		MODELS(mediaPeople).
		RETURNING(mediaPerson.AllColumns)

	util.DebugCheck(p.env, statement)

	var createdModels []model.MediaPerson
	if err := statement.QueryContext(p.ctx, p.db, &createdModels); err != nil {
		return nil, errs.BuildError(err, "could not insert media person models")
	}

	return createdModels, nil
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
	if err := statement.QueryContext(p.ctx, p.db, &createdModels); err != nil {
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

func New(env *environment.EnvironmentVariables, db *sql.DB, context context.Context) PersonRepository {
	if personRepositoryInstance != nil {
		return personRepositoryInstance
	}

	personRepositoryInstance = &personRepository{
		env:    env,
		db:     db,
		logger: logger.New(env),
		ctx:    context,
	}

	personRepositoryInstance.logger.Info("PersonRepository instance created")

	return personRepositoryInstance
}
