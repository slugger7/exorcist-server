package mediaRepository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
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
	UpdateChecksum(m model.Media) error
	GetAll(models.MediaSearchDTO) (*models.Page[models.MediaOverviewModel], error)
	GetByLibraryPathId(id uuid.UUID) ([]model.Media, error)
	GetById(id uuid.UUID) (*model.Media, error)
	Relate(model.MediaRelation) (*model.MediaRelation, error)
}

type MediaRepository struct {
	db     *sql.DB
	Env    *environment.EnvironmentVariables
	logger logger.ILogger
	ctx    context.Context
}

var mediaRepositoryInstance *MediaRepository

func New(db *sql.DB, env *environment.EnvironmentVariables, context context.Context) IMediaRepository {
	if mediaRepositoryInstance != nil {
		return mediaRepositoryInstance
	}

	mediaRepositoryInstance = &MediaRepository{
		db:     db,
		Env:    env,
		ctx:    context,
		logger: logger.New(env),
	}

	return mediaRepositoryInstance
}

func (r *MediaRepository) Create(ms []model.Media) ([]model.Media, error) {
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

	util.DebugCheck(r.Env, statement)

	models := []model.Media{}

	if err := statement.Query(r.db, &models); err != nil {
		return nil, errs.BuildError(err, "could not insert media models")
	}

	return models, nil
}

func (r *MediaRepository) UpdateExists(m model.Media) error {
	m.Modified = time.Now()

	statement := media.UPDATE().
		SET(
			media.Exists.SET(postgres.Bool(m.Exists)),
			media.Modified.SET(postgres.TimestampT(m.Modified)),
		).
		MODEL(m).
		WHERE(media.ID.EQ(postgres.UUID(m.ID)))

	util.DebugCheck(r.Env, statement)

	if _, err := statement.Exec(r.db); err != nil {
		return errs.BuildError(err, "could not update media exists: %v", m.ID)
	}

	return nil
}

func (r *MediaRepository) UpdateChecksum(m model.Media) error {
	m.Modified = time.Now()
	statement := media.UPDATE().
		SET(
			media.Checksum.SET(postgres.String(*m.Checksum)),
			media.Modified.SET(postgres.TimestampT(m.Modified)),
		).
		MODEL(m).
		WHERE(media.ID.EQ(postgres.UUID(m.ID)))

	util.DebugCheck(r.Env, statement)

	if _, err := statement.Exec(r.db); err != nil {
		return errs.BuildError(err, "could not update checksum for video: %v", m.ID)
	}

	return nil
}

func (r *MediaRepository) GetAll(search models.MediaSearchDTO) (*models.Page[models.MediaOverviewModel], error) {
	mediaRelation := table.MediaRelation
	thumbnail := table.Media.AS("thumbnail")
	selectStatement := media.SELECT(
		media.ID,
		media.Title,
		media.MediaType,
		thumbnail.ID,
	).
		FROM(
			media.LEFT_JOIN(
				mediaRelation, media.ID.EQ(mediaRelation.MediaID).
					AND(mediaRelation.RelationType.EQ(postgres.NewEnumValue(model.MediaRelationTypeEnum_Thumbnail.String()))),
			).LEFT_JOIN(
				thumbnail,
				thumbnail.ID.EQ(mediaRelation.RelatedTo),
			)).
		LIMIT(int64(search.Limit)).
		OFFSET(int64(search.Skip))

	selectStatement = helpers.OrderByDirectionColumn(search.Asc, search.OrderBy.ToColumn(), selectStatement)
	countStatement := media.SELECT(postgres.COUNT(media.ID).AS("total")).FROM(media)

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

	util.DebugCheck(r.Env, countStatement)
	util.DebugCheck(r.Env, selectStatement)

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

	return &models.Page[models.MediaOverviewModel]{
		Data:  mediaResult,
		Limit: search.Limit,
		Skip:  search.Skip,
		Total: total.Total,
	}, nil
}

func (r *MediaRepository) GetByLibraryPathId(id uuid.UUID) ([]model.Media, error) {
	statement := media.SELECT(media.Path, media.ID).
		FROM(media).
		WHERE(media.LibraryPathID.EQ(postgres.UUID(id)).
			AND(media.Exists.IS_TRUE()))

	util.DebugCheck(r.Env, statement)

	var results []model.Media
	if err := statement.QueryContext(r.ctx, r.db, &results); err != nil {
		return nil, errs.BuildError(err, "could not get media by library id: %v", id)
	}

	return results, nil
}

func (r *MediaRepository) GetById(id uuid.UUID) (*model.Media, error) {
	statement := media.SELECT(media.AllColumns).
		FROM(media).
		WHERE(media.ID.EQ(postgres.UUID(id))).
		LIMIT(1)

	util.DebugCheck(r.Env, statement)

	var result model.Media
	if err := statement.QueryContext(r.ctx, r.db, &result); err != nil {
		return nil, errs.BuildError(err, "could not get media by id: %v", id)
	}

	return &result, nil
}

func (r *MediaRepository) Relate(m model.MediaRelation) (*model.MediaRelation, error) {
	// TODO: add constraint on unique combination of media ids
	relation := table.MediaRelation
	statement := relation.INSERT(
		relation.MediaID,
		relation.RelatedTo,
		relation.RelationType,
	).
		MODEL(m).
		RETURNING(relation.AllColumns)

	util.DebugCheck(r.Env, statement)

	if err := statement.Query(r.db, &m); err != nil {
		return nil, errs.BuildError(err, "could not insert media models")
	}

	return &m, nil
}
