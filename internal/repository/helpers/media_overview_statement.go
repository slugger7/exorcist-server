package helpers

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
	"github.com/slugger7/exorcist/internal/models"
	"github.com/slugger7/exorcist/internal/repository/util"
)

type RelationFn func(relationTable postgres.ReadableTable) postgres.ReadableTable
type WhereFn func(currentWhere postgres.BoolExpression) postgres.BoolExpression

func mediaOverviewStatement(userId uuid.UUID, search dto.MediaSearchDTO, relationFn RelationFn, whereFn WhereFn) postgres.Statement {
	tagFilter := len(search.Tags) > 0
	personFilter := len(search.People) > 0
	media := table.Media
	mediaRelation := table.MediaRelation
	thumbnail := table.Media.AS("thumbnail")
	tag := table.Tag

	fromStmnt := relationFn(
		media.LEFT_JOIN(
			mediaRelation, media.ID.EQ(mediaRelation.MediaID).
				AND(mediaRelation.RelationType.EQ(postgres.NewEnumValue(model.MediaRelationTypeEnum_Thumbnail.String()))),
		).LEFT_JOIN(
			thumbnail,
			thumbnail.ID.EQ(mediaRelation.RelatedTo),
		).LEFT_JOIN(
			table.Video,
			table.Video.MediaID.EQ(media.ID),
		).LEFT_JOIN(
			table.MediaProgress,
			table.MediaProgress.MediaID.EQ(media.ID).
				AND(table.MediaProgress.UserID.EQ(postgres.UUID(userId))),
		).LEFT_JOIN(
			table.FavouriteMedia,
			table.FavouriteMedia.MediaID.EQ(media.ID).
				AND(table.FavouriteMedia.UserID.EQ(postgres.UUID(userId))),
		))

	if tagFilter {
		applyJoin := func(tbl postgres.ReadableTable) postgres.ReadableTable {
			return tbl.LEFT_JOIN(
				table.MediaTag,
				media.ID.EQ(table.MediaTag.MediaID),
			).LEFT_JOIN(
				tag,
				table.MediaTag.TagID.EQ(table.Tag.ID),
			)
		}

		fromStmnt = applyJoin(fromStmnt)
	}

	if personFilter {
		applyJoin := func(tbl postgres.ReadableTable) postgres.ReadableTable {
			return tbl.LEFT_JOIN(
				table.MediaPerson,
				media.ID.EQ(table.MediaPerson.MediaID),
			).LEFT_JOIN(
				table.Person,
				table.MediaPerson.PersonID.EQ(table.Person.ID),
			)
		}

		fromStmnt = applyJoin(fromStmnt)
	}

	selectStatement := media.SELECT(
		media.ID,
		media.Title,
		thumbnail.ID,
		table.MediaProgress.Timestamp,
		table.Video.Runtime,
		table.FavouriteMedia.ID,
		postgres.COUNT(postgres.STAR).OVER().AS("total"),
	).
		FROM(fromStmnt)

	selectStatement = OrderByDirectionColumn(search.Asc, search.OrderBy.ToColumn(), selectStatement)

	whr := media.MediaType.EQ(postgres.NewEnumValue(model.MediaTypeEnum_Primary.String())).
		AND(media.Deleted.IS_FALSE()).
		AND(media.Exists.IS_TRUE())

	if search.Favourites {
		whr = whr.AND(table.FavouriteMedia.UserID.IS_NOT_NULL())
	}

	if search.Search != "" {
		caseInsensitive := strings.ToLower(search.Search)
		likeExpression := fmt.Sprintf("%%%v%%", caseInsensitive)
		whr = whr.
			AND(
				postgres.LOWER(media.Title).LIKE(postgres.String(likeExpression)).
					OR(postgres.LOWER(media.Path).LIKE(postgres.String(likeExpression))),
			)
	}

	if tagFilter {
		tagExpressions := make([]postgres.Expression, len(search.Tags))
		for i, t := range search.Tags {
			tagExpressions[i] = postgres.String(t)
		}

		whr = whr.AND(
			tag.Name.IN(tagExpressions...),
		)
	}

	if personFilter {
		peopleExpression := make([]postgres.Expression, len(search.People))
		for i, p := range search.People {
			peopleExpression[i] = postgres.String(p)
		}

		whr = whr.AND(
			table.Person.Name.IN(peopleExpression...),
		)
	}

	if len(search.WatchStatuses) > 0 {
		var watchWhere postgres.BoolExpression
		for _, w := range search.WatchStatuses {
			var t postgres.BoolExpression
			switch w {
			case dto.WatchStatus_Watched:
				t = table.MediaProgress.Timestamp.GT(table.Video.Runtime.MUL(postgres.Float(0.9)))
			case dto.WatchStatus_Unwatched:
				t = table.MediaProgress.Timestamp.LT(table.Video.Runtime.MUL(postgres.Float(0.1))).OR(table.MediaProgress.Timestamp.IS_NULL())
			case dto.WatchStatus_InProgress:
				t = table.MediaProgress.Timestamp.GT(table.Video.Runtime.MUL(postgres.Float(0.1))).
					AND(table.MediaProgress.Timestamp.LT(table.Video.Runtime.MUL(postgres.Float(0.9))))
			default:
				continue
			}

			if watchWhere == nil {
				watchWhere = t
			} else {
				watchWhere = watchWhere.OR(t)
			}
		}

		whr = whr.AND(watchWhere)
	}

	selectStatement = selectStatement.WHERE(whereFn(whr))

	if tagFilter || personFilter {
		selectStatement = selectStatement.GROUP_BY(media.ID, thumbnail.ID, table.Video.Runtime, table.MediaProgress.Timestamp)
	}

	if tagFilter {
		selectStatement = selectStatement.HAVING(postgres.COUNT(postgres.DISTINCT(tag.Name)).EQ(postgres.Int32(int32(len(search.Tags)))))
	}

	if personFilter {
		selectStatement = selectStatement.HAVING(postgres.COUNT(postgres.DISTINCT(table.Person.Name)).EQ(postgres.Int32(int32(len(search.People)))))
	}

	selectStatement = selectStatement.
		LIMIT(int64(search.Limit)).
		OFFSET(int64(search.Skip))

	return selectStatement
}

func QueryMediaOverview(userId uuid.UUID, search dto.MediaSearchDTO, relationFn RelationFn, whereFn WhereFn, ctx context.Context, db *sql.DB, env *environment.EnvironmentVariables) (*dto.PageDTO[models.MediaOverviewModel], error) {
	selectStatement := mediaOverviewStatement(userId, search, relationFn, whereFn)

	util.DebugCheck(env, selectStatement)

	var mediaResult []struct {
		Total int
		models.MediaOverviewModel
	}
	if err := selectStatement.QueryContext(ctx, db, &mediaResult); err != nil {
		return nil, errs.BuildError(err, "could not query media for overview")
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
