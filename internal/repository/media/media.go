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

type MediaRepository interface {
	Create([]model.Media) ([]model.Media, error)
	UpdateExists(model.Media) error
	UpdateChecksum(m models.Media) error
	GetAll(userId uuid.UUID, search dto.MediaSearchDTO) (*dto.PageDTO[models.MediaOverviewModel], error)
	GetByLibraryPathId(id uuid.UUID) ([]model.Media, error)
	GetByLibraryId(libraryId uuid.UUID, pageRequest *dto.PageRequestDTO, columns postgres.ColumnList) (*dto.PageDTO[model.Media], error)
	GetById(id uuid.UUID) (*models.Media, error)
	GetByIdAndUserId(id, userId uuid.UUID) (*models.Media, error)
	Relate(model.MediaRelation) (*model.MediaRelation, error)
	Delete(m model.Media) error
	GetAssetsFor(id uuid.UUID) ([]model.Media, error)
	GetProgressForUser(id, userId uuid.UUID) (*model.MediaProgress, error)
	UpsertProgress(prog model.MediaProgress) (*model.MediaProgress, error)
	Update(m model.Media, columns postgres.ColumnList) (*model.Media, error)
}

type mediaRepository struct {
	db     *sql.DB
	env    *environment.EnvironmentVariables
	logger logger.ILogger
	ctx    context.Context
}

// GetByLibraryId implements MediaRepository.
func (r *mediaRepository) GetByLibraryId(libraryId uuid.UUID, pageRequest *dto.PageRequestDTO, columns postgres.ColumnList) (*dto.PageDTO[model.Media], error) {
	if len(columns) == 0 {
		columns = media.AllColumns
	}
	statement := media.SELECT(
		columns,
		postgres.COUNT(postgres.STAR).OVER().AS("total"),
	).
		FROM(media.
			INNER_JOIN(
				table.LibraryPath,
				media.LibraryPathID.EQ(table.LibraryPath.ID),
			),
		).
		WHERE(table.LibraryPath.LibraryID.EQ(postgres.UUID(libraryId)).
			AND(media.Deleted.IS_FALSE()).
			AND(media.Exists.IS_TRUE()))

	limit := -1
	skip := -1
	if pageRequest != nil {
		limit = pageRequest.Limit
		skip = pageRequest.Skip
		statement = statement.LIMIT(int64(pageRequest.Limit)).
			OFFSET(int64(pageRequest.Skip))
	}

	util.DebugCheck(r.env, statement)

	var mediaResult []struct {
		Total int
		model.Media
	}
	if err := statement.QueryContext(r.ctx, r.db, &mediaResult); err != nil {
		return nil, errs.BuildError(err, "querying media entities for library: %v", libraryId.String())
	}

	var data []model.Media
	total := 0
	if mediaResult != nil && len(mediaResult) > 0 {
		data = make([]model.Media, len(mediaResult))
		total = mediaResult[0].Total
		for i, o := range mediaResult {
			data[i] = o.Media
		}
	}

	return &dto.PageDTO[model.Media]{
		Data:  data,
		Limit: limit,
		Skip:  skip,
		Total: total,
	}, nil
}

// Update implements MediaRepository
func (r *mediaRepository) Update(m model.Media, columns postgres.ColumnList) (*model.Media, error) {
	if len(columns) == 0 {
		return nil, nil
	}
	columns = append(columns, media.Modified)

	m.Modified = time.Now()
	statement := media.UPDATE(columns).
		MODEL(m).
		WHERE(media.ID.EQ(postgres.UUID(m.ID))).
		RETURNING(media.ID, columns)

	var updatedModel model.Media
	if err := statement.QueryContext(r.ctx, r.db, &updatedModel); err != nil {
		return nil, errs.BuildError(err, "could not update media: %v", m.ID)
	}

	return &updatedModel, nil
}

// UpsertProgress implements MediaRepository.
func (r *mediaRepository) UpsertProgress(prog model.MediaProgress) (*model.MediaProgress, error) {
	prog.Modified = time.Now()
	mediaProgres := table.MediaProgress
	insertStatement := mediaProgres.INSERT(mediaProgres.MediaID, mediaProgres.UserID, mediaProgres.Timestamp, mediaProgres.Modified).
		MODEL(prog).
		ON_CONFLICT(mediaProgres.MediaID, mediaProgres.UserID).
		DO_UPDATE(postgres.SET(
			mediaProgres.Timestamp.SET(postgres.Float(prog.Timestamp)),
			mediaProgres.Modified.SET(postgres.LOCALTIMESTAMP()),
		)).
		RETURNING(mediaProgres.AllColumns)

	var updatedProg struct {
		model.MediaProgress
	}
	if err := insertStatement.QueryContext(r.ctx, r.db, &updatedProg); err != nil {
		return nil, errs.BuildError(err, "could not insert/update progress for user %v for media %v", prog.UserID.String(), prog.MediaID.String())
	}

	return &updatedProg.MediaProgress, nil
}

// GetProgressForUser implements MediaRepository.
func (r *mediaRepository) GetProgressForUser(id uuid.UUID, userId uuid.UUID) (*model.MediaProgress, error) {
	mediaProgress := table.MediaProgress
	selectStatement := mediaProgress.SELECT(
		mediaProgress.ID,
		mediaProgress.MediaID,
		mediaProgress.UserID,
		mediaProgress.Timestamp).
		WHERE(mediaProgress.MediaID.EQ(postgres.UUID(id)).
			AND(mediaProgress.UserID.EQ(postgres.UUID(userId))))

	var prog []struct {
		model.MediaProgress
	}
	if err := selectStatement.QueryContext(r.ctx, r.db, &prog); err != nil {
		return nil, errs.BuildError(err, "could not fetch pogress for user %v and media %v", userId.String(), id.String())
	}

	if len(prog) == 0 {
		return nil, nil
	}

	return &prog[0].MediaProgress, nil
}

// GetAssetsFor implements MediaRepository.
func (r *mediaRepository) GetAssetsFor(id uuid.UUID) ([]model.Media, error) {
	statement := media.SELECT(media.AllColumns).
		FROM(media.INNER_JOIN(table.MediaRelation, media.ID.EQ(table.MediaRelation.MediaID))).
		WHERE(table.Media.MediaType.EQ(postgres.NewEnumValue(model.MediaTypeEnum_Asset.String())))

	var entities []model.Media
	if err := statement.QueryContext(r.ctx, r.db, &entities); err != nil {
		return nil, errs.BuildError(err, "could not fetch related media for: %v", id.String())
	}

	return entities, nil
}

// Delete implements MediaRepository.
func (r *mediaRepository) Delete(m model.Media) error {
	m.Modified = time.Now()

	updateStatement := media.UPDATE(media.Exists, media.Deleted, media.Modified).
		MODEL(m).
		WHERE(media.ID.EQ(postgres.UUID(m.ID)))

	_, err := updateStatement.ExecContext(r.ctx, r.db)

	return err
}

var mediaRepositoryInstance *mediaRepository

func New(db *sql.DB, env *environment.EnvironmentVariables, context context.Context) MediaRepository {
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
	m.Media.Modified = time.Now()
	statement := media.UPDATE().
		SET(
			media.Checksum.SET(postgres.String(*m.Checksum)),
			media.Modified.SET(postgres.TimestampT(m.Media.Modified)),
		).
		MODEL(m).
		WHERE(media.ID.EQ(postgres.UUID(m.Media.ID)))

	util.DebugCheck(r.env, statement)

	if _, err := statement.Exec(r.db); err != nil {
		return errs.BuildError(err, "could not update checksum for video: %v", m.Media.ID)
	}

	return nil
}

func (r *mediaRepository) GetAll(userId uuid.UUID, search dto.MediaSearchDTO) (*dto.PageDTO[models.MediaOverviewModel], error) {
	relationFn := func(relationTable postgres.ReadableTable) postgres.ReadableTable {
		return relationTable
	}

	whereFn := func(whr postgres.BoolExpression) postgres.BoolExpression {
		return whr
	}

	mediaPage, err := helpers.QueryMediaOverview(userId, search, relationFn, whereFn, r.ctx, r.db, r.env)
	if err != nil {
		return nil, errs.BuildError(err, "colud not query media overview from media repo")
	}

	return mediaPage, nil
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
	return r.GetByIdAndUserId(id, uuid.New())
}

func (r *mediaRepository) GetByIdAndUserId(id, userId uuid.UUID) (*models.Media, error) {
	image := table.Image
	video := table.Video
	thumbnail := table.Media.AS("thumbnail")
	chapter := table.Media.AS("chapter")
	mediaChapter := table.MediaRelation.AS("media_chapter")
	mediaRelation := table.MediaRelation
	mediaPerson := table.MediaPerson
	person := table.Person
	mediaTag := table.MediaTag
	tag := table.Tag

	statement := media.SELECT(
		media.AllColumns,
		image.AllColumns,
		video.AllColumns,
		thumbnail.ID,
		person.AllColumns,
		tag.AllColumns,
		table.MediaProgress.Timestamp,
		table.FavouriteMedia.ID,
		mediaChapter.Metadata,
		mediaChapter.RelatedTo,
	).FROM(media.
		LEFT_JOIN(image, image.MediaID.EQ(media.ID)).
		LEFT_JOIN(video, video.MediaID.EQ(media.ID)).
		LEFT_JOIN(mediaRelation, mediaRelation.MediaID.EQ(media.ID).
			AND(mediaRelation.RelationType.EQ(
				postgres.NewEnumValue(model.MediaRelationTypeEnum_Thumbnail.String()),
			))).
		LEFT_JOIN(thumbnail, thumbnail.ID.EQ(mediaRelation.RelatedTo).
			AND(thumbnail.MediaType.EQ(postgres.NewEnumValue(model.MediaTypeEnum_Asset.String())))).
		LEFT_JOIN(mediaChapter, mediaChapter.MediaID.EQ(media.ID).
			AND(mediaChapter.RelationType.EQ(postgres.NewEnumValue(model.MediaRelationTypeEnum_Chapter.String())))).
		LEFT_JOIN(chapter, chapter.ID.EQ(mediaChapter.RelatedTo).
			AND(chapter.MediaType.EQ(postgres.NewEnumValue(model.MediaTypeEnum_Asset.String())))).
		LEFT_JOIN(mediaPerson, mediaPerson.MediaID.EQ(media.ID)).
		LEFT_JOIN(person, person.ID.EQ(mediaPerson.PersonID)).
		LEFT_JOIN(mediaTag, mediaTag.MediaID.EQ(media.ID)).
		LEFT_JOIN(tag, tag.ID.EQ(mediaTag.TagID)).
		LEFT_JOIN(table.MediaProgress, table.MediaProgress.MediaID.EQ(media.ID).AND(table.MediaProgress.UserID.EQ(postgres.UUID(userId)))).
		LEFT_JOIN(table.FavouriteMedia, table.FavouriteMedia.MediaID.EQ(media.ID).AND(table.FavouriteMedia.UserID.EQ(postgres.UUID(userId)))),
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
		relation.Metadata,
	).
		MODEL(m).
		RETURNING(relation.AllColumns)

	util.DebugCheck(r.env, statement)

	if err := statement.Query(r.db, &m); err != nil {
		return nil, errs.BuildError(err, "could not insert media models")
	}

	return &m, nil
}
