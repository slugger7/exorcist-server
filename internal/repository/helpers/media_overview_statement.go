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
				))).
		LIMIT(int64(search.Limit)).
		OFFSET(int64(search.Skip))

	selectStatement = OrderByDirectionColumn(search.Asc, search.OrderBy.ToColumn(), selectStatement)
	countStatement := media.SELECT(postgres.COUNT(media.ID).AS("total")).FROM(relationFn(media))

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

	return selectStatement, countStatement
}
