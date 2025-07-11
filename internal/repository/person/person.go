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
	GetAll() ([]model.Person, error)
	GetMedia(id uuid.UUID, search dto.MediaSearchDTO) (*dto.PageDTO[models.MediaOverviewModel], error)
}

type personRepository struct {
	env    *environment.EnvironmentVariables
	db     *sql.DB
	logger logger.ILogger
	ctx    context.Context
}

func innerJoinMediaPerson(id uuid.UUID, t postgres.ReadableTable) postgres.ReadableTable {
	media := table.Media
	mediaPerson := table.MediaPerson
	return t.INNER_JOIN(
		mediaPerson,
		media.ID.EQ(mediaPerson.MediaID).
			AND(mediaPerson.PersonID.EQ(postgres.UUID(id))),
	)
}

// GetMedia implements PersonRepository.
func (r *personRepository) GetMedia(id uuid.UUID, search dto.MediaSearchDTO) (*dto.PageDTO[models.MediaOverviewModel], error) {
	media := table.Media
	mediaRelation := table.MediaRelation
	thumbnail := table.Media.AS("thumbnail")

	selectStatement := media.SELECT(
		media.ID,
		media.Title,
		media.MediaType,
		thumbnail.ID,
	).
		FROM(
			innerJoinMediaPerson(id, media.LEFT_JOIN(
				mediaRelation, media.ID.EQ(mediaRelation.MediaID).
					AND(mediaRelation.RelationType.EQ(postgres.NewEnumValue(model.MediaRelationTypeEnum_Thumbnail.String()))),
			).LEFT_JOIN(
				thumbnail,
				thumbnail.ID.EQ(mediaRelation.RelatedTo),
			))).
		LIMIT(int64(search.Limit)).
		OFFSET(int64(search.Skip))

	selectStatement = helpers.OrderByDirectionColumn(search.Asc, search.OrderBy.ToColumn(), selectStatement)
	countStatement := media.SELECT(postgres.COUNT(media.ID).AS("total")).FROM(innerJoinMediaPerson(id, media))

	whr := media.MediaType.EQ(postgres.NewEnumValue(model.MediaTypeEnum_Primary.String())).
		AND(media.Deleted.IS_FALSE()).
		AND(media.Exists.IS_TRUE())

	if search.Search != "" {
		caseInsensitive := strings.ToLower(search.Search)
		likeExpression := fmt.Sprintf("%%%v%%", caseInsensitive)
		whr = whr.
			AND(
				postgres.LOWER(media.Title).LIKE(postgres.String(likeExpression)).
					OR(postgres.LOWER(media.Path).LIKE(postgres.String(likeExpression))),
			)
	}

	selectStatement = selectStatement.WHERE(whr)
	countStatement = countStatement.WHERE(whr)

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
func (p *personRepository) GetAll() ([]model.Person, error) {
	statement := person.SELECT(person.AllColumns)

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
