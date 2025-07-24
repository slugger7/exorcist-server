package helpers

import (
	"fmt"
	"strings"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/dto"
)

func MediaOverviewStatement(search dto.MediaSearchDTO, relationFn func(relationTable postgres.ReadableTable) postgres.ReadableTable) (postgres.Statement, postgres.Statement) {
	media := table.Media
	mediaRelation := table.MediaRelation
	thumbnail := table.Media.AS("thumbnail")
	tag := table.Tag.AS("tag")
	selectStatement := media.SELECT(
		media.ID,
		media.Title,
		media.MediaType,
		thumbnail.ID,
	).
		FROM(
			relationFn(
				media.LEFT_JOIN(
					mediaRelation, media.ID.EQ(mediaRelation.MediaID).
						AND(mediaRelation.RelationType.EQ(postgres.NewEnumValue(model.MediaRelationTypeEnum_Thumbnail.String()))),
				).LEFT_JOIN(
					thumbnail,
					thumbnail.ID.EQ(mediaRelation.RelatedTo),
				).LEFT_JOIN(
					table.MediaTag,
					media.ID.EQ(table.MediaTag.MediaID),
				).LEFT_JOIN(
					tag,
					table.MediaTag.TagID.EQ(table.Tag.ID),
				)))

	selectStatement = OrderByDirectionColumn(search.Asc, search.OrderBy.ToColumn(), selectStatement)
	countStatement := media.SELECT(postgres.COUNT(media.ID).AS("total")).FROM(relationFn(media).LEFT_JOIN(
		table.MediaTag,
		media.ID.EQ(table.MediaTag.MediaID),
	).LEFT_JOIN(
		tag,
		table.MediaTag.TagID.EQ(table.Tag.ID),
	))

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

	if len(search.Tags) > 0 {
		tagExpressions := make([]postgres.Expression, len(search.Tags))
		for i, t := range search.Tags {
			tagExpressions[i] = postgres.String(t)
		}

		whr = whr.AND(
			tag.Name.IN(tagExpressions...),
		)
	}

	selectStatement = selectStatement.WHERE(whr)
	countStatement = countStatement.WHERE(whr)

	if len(search.Tags) > 0 {
		selectStatement = selectStatement.GROUP_BY(media.ID, thumbnail.ID).
			HAVING(postgres.COUNT(postgres.DISTINCT(tag.Name)).EQ(postgres.Int32(int32(len(search.Tags)))))
		countStatement = countStatement.GROUP_BY(media.ID).
			HAVING(postgres.COUNT(postgres.DISTINCT(tag.Name)).EQ(postgres.Int32(int32(len(search.Tags)))))
	}

	selectStatement = selectStatement.
		LIMIT(int64(search.Limit)).
		OFFSET(int64(search.Skip))

	return selectStatement, countStatement
}
