package media

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/models"
	"github.com/slugger7/exorcist/internal/repository/helpers"
	"github.com/slugger7/exorcist/internal/repository/util"
)

type IMediaRepository interface {
	Create([]model.Media) ([]model.Media, error)
	UpdateExists(model.Media) error
	GetAll(models.MediaSearchDTO) (*models.Page[models.MediaOverviewModel], error)
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
	media := table.Media

	statement := media.INSERT(
		media.LibraryPathID,
		media.Path,
		media.FileName,
		media.Title,
		media.Size,
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

	statement := table.Video.UPDATE().
		SET(
			table.Media.Exists.SET(postgres.Bool(m.Exists)),
			table.Media.Modified.SET(postgres.TimestampT(m.Modified)),
		).
		MODEL(m).
		WHERE(table.Video.ID.EQ(postgres.UUID(m.ID)))

	util.DebugCheck(r.Env, statement)

	if _, err := statement.Exec(r.db); err != nil {
		return errs.BuildError(err, "could not update media exists: %v", m.ID)
	}

	return nil
}

func (r *MediaRepository) UpdateChecksum(m model.Media) error {
	m.Modified = time.Now()
	statement := table.Video.UPDATE().
		SET(
			table.Media.Checksum.SET(postgres.String(*m.Checksum)),
			table.Media.Modified.SET(postgres.TimestampT(m.Modified)),
		).
		MODEL(m).
		WHERE(table.Video.ID.EQ(postgres.UUID(m.ID)))

	util.DebugCheck(r.Env, statement)

	if _, err := statement.Exec(r.db); err != nil {
		return errs.BuildError(err, "could not update checksum for video: %v", m.ID)
	}

	return nil
}

func (r *MediaRepository) GetAll(search models.MediaSearchDTO) (*models.Page[models.MediaOverviewModel], error) {
	media := table.Media
	mediaRelation := table.MediaRelation
	thumbnail := table.Media.AS("thumbnail")
	image := table.Image
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
			).LEFT_JOIN(
				image,
				image.MediaID.EQ(thumbnail.ID),
			)).
		LIMIT(int64(search.Limit)).
		OFFSET(int64(search.Skip))

	selectStatement = helpers.OrderByDirectionColumn(search.Asc, search.OrderBy.ToColumn(), selectStatement)
	countStatement := media.SELECT(postgres.COUNT(media.ID).AS("total")).FROM(media)

	if search.Search != "" {
		caseInsensitive := strings.ToLower(search.Search)
		likeExpression := fmt.Sprintf("%%%v%%", caseInsensitive)
		query := media.Deleted.IS_FALSE().
			AND(media.Exists.IS_TRUE()).
			AND(
				postgres.LOWER(media.Title).LIKE(postgres.String(likeExpression)).
					OR(postgres.LOWER(media.Path).LIKE(postgres.String(likeExpression))),
			)

		selectStatement = selectStatement.WHERE(query)
		countStatement = countStatement.WHERE(query)
	}

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
