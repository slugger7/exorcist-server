package tagRepository

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

var tag = table.Tag
var mediaTag = table.MediaTag

type TagRepository interface {
	GetByName(name string) (*model.Tag, error)
	Create(names []string) ([]model.Tag, error)
	AddToMedia(mediaPeople []model.MediaTag) ([]model.MediaTag, error)
	RemoveFromMedia(mediaTag model.MediaTag) error
}

type tagRepository struct {
	env    *environment.EnvironmentVariables
	db     *sql.DB
	logger logger.ILogger
	ctx    context.Context
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
func (p *tagRepository) AddToMedia(mediaPeople []model.MediaTag) ([]model.MediaTag, error) {
	if len(mediaPeople) == 0 {
		return nil, nil
	}

	statement := mediaTag.INSERT(mediaTag.MediaID, mediaTag.TagID).
		MODELS(mediaPeople).
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
