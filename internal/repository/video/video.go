package videoRepository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/models"
	"github.com/slugger7/exorcist/internal/repository/helpers"
)

type VideoLibraryPathModel struct {
	model.Video
	model.LibraryPath
}

type IVideoRepository interface {
	GetAll() ([]model.Video, error)
	GetByLibraryPathId(id uuid.UUID) ([]model.Video, error)
	GetById(id uuid.UUID) (*model.Video, error)
	UpdateExists(video *model.Video) error
	UpdateChecksum(video *model.Video) error
	Insert(models []model.Video) ([]model.Video, error)
	GetByIdWithLibraryPath(id uuid.UUID) (*VideoLibraryPathModel, error)
	GetOverview(limit, skip int, ordinal *models.VideoOrdinal, asc bool) (*models.Page[models.VideoOverviewModel], error)
}

type VideoRepository struct {
	db  *sql.DB
	Env *environment.EnvironmentVariables
}

var videoRepoInstance *VideoRepository

func New(db *sql.DB, env *environment.EnvironmentVariables) IVideoRepository {
	if videoRepoInstance != nil {
		return videoRepoInstance
	}
	videoRepoInstance = &VideoRepository{
		db:  db,
		Env: env,
	}
	return videoRepoInstance
}

func (vr *VideoRepository) GetAll() ([]model.Video, error) {
	statement := table.Video.SELECT(table.Video.AllColumns).
		FROM(table.Video)

	var vids []struct{ model.Video }
	if err := statement.Query(vr.db, &vids); err != nil {
		return nil, errs.BuildError(err, "could not get all videos")
	}

	if vids == nil {
		return nil, nil
	}

	var models []model.Video
	for _, v := range vids {
		models = append(models, v.Video)
	}

	return models, nil
}

func (ds *VideoRepository) GetByLibraryPathId(id uuid.UUID) ([]model.Video, error) {
	var vids []struct{ model.Video }
	if err := ds.getByLibraryPathIdStatement(id).Query(&vids); err != nil {
		return nil, errs.BuildError(err, "could not get videos by library path id: %v", id)
	}

	vidModels := []model.Video{}
	for _, v := range vids {
		vidModels = append(vidModels, v.Video)
	}

	return vidModels, nil
}

func (i *VideoRepository) UpdateExists(v *model.Video) error {
	v.Modified = time.Now()
	_, err := i.updateVideoExistsStatement(*v).Exec()
	if err != nil {
		return errs.BuildError(err, "could not update video exists: %v", v.ID)
	}
	return nil
}

func (ds *VideoRepository) Insert(models []model.Video) ([]model.Video, error) {
	if len(models) == 0 {
		return nil, nil
	}

	var vids []struct{ model.Video }

	if err := ds.insertStatement(models).Query(&vids); err != nil {
		return nil, errs.BuildError(err, "could not insert video models to database")
	}

	var vidModels = []model.Video{}
	for _, v := range vids {
		vidModels = append(vidModels, v.Video)
	}

	return vidModels, nil
}

func (ds *VideoRepository) GetById(id uuid.UUID) (*model.Video, error) {
	var vids []struct{ model.Video }
	if err := ds.getByIdStatement(id).Query(&vids); err != nil {
		return nil, errs.BuildError(err, "error getting video from db for id %v", id)
	}

	var video *model.Video
	if len(vids) == 1 {
		video = &vids[len(vids)-1].Video
	}

	return video, nil
}

func (ds *VideoRepository) GetByIdWithLibraryPath(id uuid.UUID) (*VideoLibraryPathModel, error) {
	var results []VideoLibraryPathModel
	if err := ds.getByIdWithLibraryPathStatement(id).Query(&results); err != nil {
		return nil, errs.BuildError(err, "error getting video by id (%v) with library path", id.String())
	}

	var result VideoLibraryPathModel
	if results != nil {
		result = results[len(results)-1]
	}

	return &result, nil
}

func (ds *VideoRepository) UpdateChecksum(video *model.Video) error {
	video.Modified = time.Now()
	if _, err := ds.updateChecksumStatement(*video).Exec(); err != nil {
		return errs.BuildError(err, "error updating video (%v) checksum %v", video.ID, video.Checksum)
	}

	return nil
}

func (ds *VideoRepository) GetOverview(limit, skip int, ordinal *models.VideoOrdinal, asc bool) (*models.Page[models.VideoOverviewModel], error) {

	selectStatement := table.Video.SELECT(
		table.Video.ID,
		table.Video.RelativePath,
		table.LibraryPath.Path,
		table.Video.Title,
		table.Image.Path,
		table.VideoImage.VideoImageType,
	).
		FROM(table.Video.
			INNER_JOIN(
				table.LibraryPath,
				table.Video.LibraryPathID.EQ(table.LibraryPath.ID)).
			LEFT_JOIN(
				table.VideoImage,
				table.VideoImage.VideoID.EQ(table.Video.ID).
					AND(table.VideoImage.VideoImageType.EQ(
						postgres.NewEnumValue(model.VideoImageTypeEnum_Thumbnail.String())))).
			INNER_JOIN(
				table.Image,
				table.Image.ID.EQ(table.VideoImage.ImageID),
			)).
		LIMIT(int64(limit)).
		OFFSET(int64(skip))

	if ordinal != nil {
		selectStatement = helpers.OrderByDirectionColumn(asc, ordinal.ToColumn(), selectStatement)
	}

	countStatement := table.Video.SELECT(postgres.COUNT(table.Video.ID).AS("total")).FROM(table.Video)

	sql := countStatement.DebugSql()
	var vids []models.VideoOverviewModel

	fmt.Println(sql)

	if err := selectStatement.Query(ds.db, &vids); err != nil {
		return nil, errs.BuildError(err, "could not query videos for overview")
	}

	var res struct {
		Total int
	}

	if err := countStatement.Query(ds.db, &res); err != nil {
		return nil, errs.BuildError(err, "could not query videos for overview total")
	}
	return &models.Page[models.VideoOverviewModel]{
		Data:  vids,
		Limit: limit,
		Skip:  skip,
		Total: res.Total,
	}, nil
}
