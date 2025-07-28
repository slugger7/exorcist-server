package tagRepository

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

var tag = table.Tag
var mediaTag = table.MediaTag

type TagRepository interface {
	GetByName(name string) (*model.Tag, error)
	Create(names []string) ([]model.Tag, error)
	AddToMedia(mediaTags []model.MediaTag) ([]model.MediaTag, error)
	RemoveFromMedia(mediaTag model.MediaTag) error
	GetAll(search dto.TagSearchDTO) ([]model.Tag, error)
	GetById(id uuid.UUID) (*model.Tag, error)
	GetMedia(id uuid.UUID, search dto.MediaSearchDTO) (*dto.PageDTO[models.MediaOverviewModel], error)
}

type tagRepository struct {
	env    *environment.EnvironmentVariables
	db     *sql.DB
	logger logger.ILogger
	ctx    context.Context
}

// GetMedia implements TagRepository.
func (r *tagRepository) GetMedia(id uuid.UUID, search dto.MediaSearchDTO) (*dto.PageDTO[models.MediaOverviewModel], error) {
	search.Tags = []string{}
	relationFn := func(relationTable postgres.ReadableTable) postgres.ReadableTable {
		media := table.Media
		mediaTag := table.MediaTag
		return relationTable.INNER_JOIN(
			mediaTag,
			media.ID.EQ(mediaTag.MediaID),
		)
	}

	whereFn := func(whr postgres.BoolExpression) postgres.BoolExpression {
		return whr.AND(mediaTag.TagID.EQ(postgres.UUID(id)))
	}

	selectStatement := helpers.MediaOverviewStatement(search, relationFn, whereFn)

	util.DebugCheck(r.env, selectStatement)

	var mediaResult []struct {
		Total int
		models.MediaOverviewModel
	}
	if err := selectStatement.QueryContext(r.ctx, r.db, &mediaResult); err != nil {
		return nil, errs.BuildError(err, "could not query media by tag id: %v", id.String())
	}

	data := make([]models.MediaOverviewModel, len(mediaResult))
	total := 0
	if mediaResult != nil && len(mediaResult) > 0 {
		total = mediaResult[0].Total
		for i, o := range mediaResult {
			data[i] = o.MediaOverviewModel
		}
	}

	return &dto.PageDTO[models.MediaOverviewModel]{
		Data:  data,
		Limit: search.Limit,
		Skip:  search.Skip,
		Total: total,
	}, nil
}

// GetById implements TagRepository.
func (p *tagRepository) GetById(id uuid.UUID) (*model.Tag, error) {
	statement := tag.SELECT(tag.AllColumns).
		WHERE(tag.ID.EQ(postgres.UUID(id)))

	var tagModels []model.Tag
	if err := statement.QueryContext(p.ctx, p.db, &tagModels); err != nil {
		return nil, errs.BuildError(err, "colud not get tag by id from db %v", id)
	}

	if len(tagModels) == 0 {
		return nil, nil
	}

	return &tagModels[0], nil
}

// GetAll implements TagRepository.
func (p *tagRepository) GetAll(search dto.TagSearchDTO) ([]model.Tag, error) {
	statement := tag.SELECT(tag.AllColumns).
		FROM(
			tag.LEFT_JOIN(table.MediaTag, tag.ID.EQ(table.MediaTag.TagID)),
		).
		GROUP_BY(tag.ID).
		ORDER_BY(search.ToOrderByClause()...)

	if search.Search != "" {
		caseInsensitive := strings.ToLower(search.Search)
		statement = statement.WHERE(postgres.LOWER(tag.Name).LIKE(postgres.String(fmt.Sprintf("%%%v%%", caseInsensitive))))
	}

	var tags []model.Tag
	if err := statement.QueryContext(p.ctx, p.db, &tags); err != nil {
		return nil, errs.BuildError(err, "could not fetch tags from database")
	}

	return tags, nil
}

// RemoveFromMedia implements ITagRepository.
func (p *tagRepository) RemoveFromMedia(mp model.MediaTag) error {
	statement := mediaTag.DELETE().
		WHERE(mediaTag.MediaID.EQ(postgres.UUID(mp.MediaID)).
			AND(mediaTag.TagID.EQ(postgres.UUID(mp.TagID))))

	util.DebugCheck(p.env, statement)

	if _, err := statement.ExecContext(p.ctx, p.db); err != nil {
		return errs.BuildError(err, "could not delete media tag with media id %v and tag id %v", mp.MediaID.String(), mp.TagID.String())
	}

	return nil
}

// AddToMedia implements ITagRepository.
func (p *tagRepository) AddToMedia(mediaTags []model.MediaTag) ([]model.MediaTag, error) {
	if len(mediaTags) == 0 {
		return nil, nil
	}

	statement := mediaTag.INSERT(mediaTag.MediaID, mediaTag.TagID).
		MODELS(mediaTags).
		RETURNING(mediaTag.AllColumns)

	util.DebugCheck(p.env, statement)

	var createdModels []model.MediaTag
	if err := statement.QueryContext(p.ctx, p.db, &createdModels); err != nil {
		return nil, errs.BuildError(err, "could not insert media tag models")
	}

	return createdModels, nil
}

// Create implements ITagRepository.
func (p *tagRepository) Create(names []string) ([]model.Tag, error) {
	if len(names) == 0 {
		return nil, nil
	}

	peopleModels := make([]model.Tag, len(names))
	for i, n := range names {
		peopleModels[i] = model.Tag{Name: n}
	}

	statement := tag.INSERT(tag.Name).
		MODELS(peopleModels).
		RETURNING(tag.AllColumns)

	util.DebugCheck(p.env, statement)

	var createdModels []model.Tag
	if err := statement.QueryContext(p.ctx, p.db, &createdModels); err != nil {
		return nil, errs.BuildError(err, "could not insert new people models")
	}

	return createdModels, nil
}

// GetByName implements ITagRepository.
func (p *tagRepository) GetByName(name string) (*model.Tag, error) {
	statement := tag.SELECT(tag.AllColumns).
		FROM(tag).
		WHERE(postgres.LOWER(tag.Name).EQ(postgres.String(strings.ToLower(name)))).
		LIMIT(1)

	util.DebugCheck(p.env, statement)

	var tagModel []model.Tag

	if err := statement.QueryContext(p.ctx, p.db, &tagModel); err != nil {
		return nil, errs.BuildError(err, "could not query people by name")
	}

	if len(tagModel) == 0 {
		return nil, nil
	}

	return &tagModel[0], nil
}

var tagRepositoryInstance *tagRepository

func New(env *environment.EnvironmentVariables, db *sql.DB, context context.Context) TagRepository {
	if tagRepositoryInstance == nil {
		tagRepositoryInstance = &tagRepository{
			env:    env,
			db:     db,
			logger: logger.New(env),
			ctx:    context,
		}

		tagRepositoryInstance.logger.Info("TagRepository instance created")
	}

	return tagRepositoryInstance
}
