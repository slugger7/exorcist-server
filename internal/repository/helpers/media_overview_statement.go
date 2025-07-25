package helpers

import (
	"fmt"
	"strings"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/dto"
)

type RelationFn func(relationTable postgres.ReadableTable) postgres.ReadableTable
type WhereFn func(currentWhere postgres.BoolExpression) postgres.BoolExpression

func MediaOverviewStatement(search dto.MediaSearchDTO, relationFn RelationFn, whereFn WhereFn) (postgres.Statement, postgres.Statement) {
	createSelect := func() postgres.SelectStatement {
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
			))

		countFromStmnt := relationFn(media)

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
			countFromStmnt = applyJoin(countFromStmnt)
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
			countFromStmnt = applyJoin(countFromStmnt)
		}

		selectStatement := media.SELECT(
			media.ID,
			media.Title,
			media.MediaType,
			thumbnail.ID,
		).
			FROM(fromStmnt)

		selectStatement = OrderByDirectionColumn(search.Asc, search.OrderBy.ToColumn(), selectStatement)

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

		selectStatement = selectStatement.WHERE(whereFn(whr))

		if tagFilter || personFilter {
			selectStatement = selectStatement.GROUP_BY(media.ID, thumbnail.ID, table.Video.Runtime)
		}

		if tagFilter {
			applyHaving := func(sl postgres.SelectStatement) postgres.SelectStatement {
				return sl.HAVING(postgres.COUNT(postgres.DISTINCT(tag.Name)).EQ(postgres.Int32(int32(len(search.Tags)))))
			}

			selectStatement = applyHaving(selectStatement)
		}

		if personFilter {
			applyHaving := func(sl postgres.SelectStatement) postgres.SelectStatement {
				return sl.HAVING(postgres.COUNT(postgres.DISTINCT(table.Person.Name)).EQ(postgres.Int32(int32(len(search.People)))))
			}

			selectStatement = applyHaving(selectStatement)
		}

		return selectStatement
	}

	countStatement := createSelect().AsTable("media_overview").SELECT(postgres.COUNT(postgres.String("id")).AS("total"))

	selectStatement := createSelect().
		LIMIT(int64(search.Limit)).
		OFFSET(int64(search.Skip))

	return selectStatement, countStatement
}
