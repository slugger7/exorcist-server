package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
)

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
