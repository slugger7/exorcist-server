package dto

import (
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
)

type PersonOrdinal string

const (
	PersonOrdinal_MediaCount PersonOrdinal = "count"
	PersonOrdinal_Name       PersonOrdinal = "name"
)

var PersonOrdinalAllValues = []PersonOrdinal{
	PersonOrdinal_MediaCount,
	PersonOrdinal_Name,
}

func (o PersonSearchDTO) ToOrderByClause() []postgres.OrderByClause {
	person := table.Person

	arr := []postgres.OrderByClause{}
	switch o.OrderBy {
	case PersonOrdinal_MediaCount:
		if o.Asc {
			arr = append(arr, postgres.COUNT(person.ID).ASC())
		} else {
			arr = append(arr, postgres.COUNT(person.ID).DESC())
		}
		arr = append(arr, person.Name)
	case PersonOrdinal_Name:
		if o.Asc {
			arr = append(arr, person.Name.ASC())
		} else {
			arr = append(arr, person.Name.DESC())
		}
	default:
		if o.Asc {
			arr = append(arr, person.Name.ASC())
		} else {
			arr = append(arr, person.Name.DESC())
		}
	}

	return arr
}

func (o PersonOrdinal) String() string {
	return string(o)
}

type PersonDTO struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}

func (o *PersonDTO) FromModel(m *model.Person) *PersonDTO {
	o.ID = m.ID
	o.Name = m.Name
	o.Created = m.Created
	o.Modified = m.Modified

	return o
}

type PersonSearchDTO struct {
	Search  string        `form:"search" json:"search"`
	OrderBy PersonOrdinal `form:"orderBy" json:"orderBy"`
	Asc     bool          `form:"asc" json:"asc"`
}

type PersonUpdateDTO struct {
	Name string `json:"name"`
}
