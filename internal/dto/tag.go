package dto

import (
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
)

type TagOrdinal string

const (
	TagOrdinal_MediaCount TagOrdinal = "count"
	TagOrdinal_Name       TagOrdinal = "name"
)

var TagOrdinalAllValues = []TagOrdinal{
	TagOrdinal_MediaCount,
	TagOrdinal_Name,
}

func (o TagSearchDTO) ToOrderByClause() []postgres.OrderByClause {
	tag := table.Tag

	arr := []postgres.OrderByClause{}
	switch o.OrderBy {
	case TagOrdinal_MediaCount:
		if o.Asc {
			arr = append(arr, postgres.COUNT(tag.ID).ASC())
		} else {
			arr = append(arr, postgres.COUNT(tag.ID).DESC())
		}
		arr = append(arr, tag.Name)
	case TagOrdinal_Name:
		if o.Asc {
			arr = append(arr, tag.Name.ASC())
		} else {
			arr = append(arr, tag.Name.DESC())
		}
	default:
		if o.Asc {
			arr = append(arr, tag.Name.ASC())
		} else {
			arr = append(arr, tag.Name.DESC())
		}
	}

	return arr
}

func (o TagOrdinal) String() string {
	return string(o)
}

type TagDTO struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}

func (o *TagDTO) FromModel(m *model.Tag) *TagDTO {
	o.ID = m.ID
	o.Name = m.Name
	o.Created = m.Created
	o.Modified = m.Modified

	return o
}

type TagSearchDTO struct {
	Search  string     `form:"search" json:"search"`
	OrderBy TagOrdinal `form:"orderBy" json:"orderBy"`
	Asc     bool       `form:"asc" json:"asc"`
}
