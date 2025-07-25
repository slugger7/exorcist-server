package mediaRepository

import (
	"context"
	"database/sql"
	"time"

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

var media = table.Media

type IMediaRepository interface {
	Create([]model.Media) ([]model.Media, error)
	UpdateExists(model.Media) error
	UpdateChecksum(m models.Media) error
	GetAll(dto.MediaSearchDTO) (*dto.PageDTO[models.MediaOverviewModel], error)
	GetByLibraryPathId(id uuid.UUID) ([]model.Media, error)
	GetById(id uuid.UUID) (*models.Media, error)
	Relate(model.MediaRelation) (*model.MediaRelation, error)
}

type mediaRepository struct {
	db     *sql.DB
	env    *environment.EnvironmentVariables
	logger logger.ILogger
	ctx    context.Context
}

var mediaRepositoryInstance *mediaRepository

func New(db *sql.DB, env *environment.EnvironmentVariables, context context.Context) IMediaRepository {
	if mediaRepositoryInstance != nil {
		return mediaRepositoryInstance
	}

	mediaRepositoryInstance = &mediaRepository{
		db:     db,
		env:    env,
		ctx:    context,
		logger: logger.New(env),
	}

	return mediaRepositoryInstance
}

func (r *mediaRepository) Create(ms []model.Media) ([]model.Media, error) {
	if len(ms) == 0 {
		return nil, nil
	}

	statement := media.INSERT(
		media.LibraryPathID,
		media.Path,
		media.Title,
		media.Size,
		media.MediaType,
	).
		MODELS(ms).
		RETURNING(media.AllColumns)

	util.DebugCheck(r.env, statement)

	models := []model.Media{}

	if err := statement.Query(r.db, &models); err != nil {
		return nil, errs.BuildError(err, "could not insert media models")
	}

	return models, nil
}

func (r *mediaRepository) UpdateExists(m model.Media) error {
	m.Modified = time.Now()

	statement := media.UPDATE().
		SET(
			media.Exists.SET(postgres.Bool(m.Exists)),
			media.Modified.SET(postgres.TimestampT(m.Modified)),
		).
		MODEL(m).
		WHERE(media.ID.EQ(postgres.UUID(m.ID)))

	util.DebugCheck(r.env, statement)

	if _, err := statement.Exec(r.db); err != nil {
		return errs.BuildError(err, "could not update media exists: %v", m.ID)
	}

	return nil
}

func (r *mediaRepository) UpdateChecksum(m models.Media) error {
	m.Modified = time.Now()
	statement := media.UPDATE().
		SET(
			media.Checksum.SET(postgres.String(*m.Checksum)),
			media.Modified.SET(postgres.TimestampT(m.Modified)),
		).
		MODEL(m).
		WHERE(media.ID.EQ(postgres.UUID(m.Media.ID)))

	util.DebugCheck(r.env, statement)

	if _, err := statement.Exec(r.db); err != nil {
		return errs.BuildError(err, "could not update checksum for video: %v", m.Media.ID)
	}

	return nil
}

func (r *mediaRepository) GetAll(search dto.MediaSearchDTO) (*dto.PageDTO[models.MediaOverviewModel], error) {
	relationFn := func(relationTable postgres.ReadableTable) postgres.ReadableTable {
		return relationTable
	}

	whereFn := func(whr postgres.BoolExpression) postgres.BoolExpression {
		return whr
	}

	selectStatement, countStatement := helpers.MediaOverviewStatement(search, relationFn, whereFn)

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

func (r *mediaRepository) GetByLibraryPathId(id uuid.UUID) ([]model.Media, error) {
	statement := media.SELECT(media.Path, media.ID).
		FROM(media).
		WHERE(media.LibraryPathID.EQ(postgres.UUID(id)).
			AND(media.Exists.IS_TRUE()))

	util.DebugCheck(r.env, statement)

	var results []model.Media
	if err := statement.QueryContext(r.ctx, r.db, &results); err != nil {
		return nil, errs.BuildError(err, "could not get media by library id: %v", id)
	}

	return results, nil
}

func (r *mediaRepository) GetById(id uuid.UUID) (*models.Media, error) {
	image := table.Image
	video := table.Video
	thumbnail := table.Media.AS("thumbnail")
	mediaRelation := table.MediaRelation
	mediaPerson := table.MediaPerson
	person := table.Person
	mediaTag := table.MediaTag
	tag := table.Tag
	statement := media.SELECT(media.AllColumns, image.AllColumns, video.AllColumns, thumbnail.ID, person.AllColumns, tag.AllColumns).
		FROM(media.
			LEFT_JOIN(image, image.MediaID.EQ(media.ID)).
			LEFT_JOIN(video, video.MediaID.EQ(media.ID)).
			LEFT_JOIN(mediaRelation, mediaRelation.MediaID.EQ(media.ID).
				AND(mediaRelation.RelationType.EQ(
					postgres.NewEnumValue(model.MediaRelationTypeEnum_Thumbnail.String()),
				))).
			LEFT_JOIN(thumbnail, thumbnail.ID.EQ(mediaRelation.RelatedTo).
				AND(thumbnail.MediaType.EQ(postgres.NewEnumValue(model.MediaTypeEnum_Asset.String())))).
			LEFT_JOIN(mediaPerson, mediaPerson.MediaID.EQ(media.ID)).
			LEFT_JOIN(person, person.ID.EQ(mediaPerson.PersonID)).
			LEFT_JOIN(mediaTag, mediaTag.MediaID.EQ(media.ID)).
			LEFT_JOIN(tag, tag.ID.EQ(mediaTag.TagID)),
		).
		WHERE(media.ID.EQ(postgres.UUID(id)))

	util.DebugCheck(r.env, statement)

	var result models.Media
	if err := statement.QueryContext(r.ctx, r.db, &result); err != nil {
		return nil, errs.BuildError(err, "could not get media by id: %v", id)
	}

	return &result, nil
}

func (r *mediaRepository) Relate(m model.MediaRelation) (*model.MediaRelation, error) {
	// TODO: add constraint on unique combination of media ids
	relation := table.MediaRelation
	statement := relation.INSERT(
		relation.MediaID,
		relation.RelatedTo,
		relation.RelationType,
	).
		MODEL(m).
		RETURNING(relation.AllColumns)

	util.DebugCheck(r.env, statement)

	if err := statement.Query(r.db, &m); err != nil {
		return nil, errs.BuildError(err, "could not insert media models")
	}

	return &m, nil
}
