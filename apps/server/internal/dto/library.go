package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
)

type CreateLibraryDTO struct {
	Name string `json:"name" binding:"required"`
}

type LibraryDTO struct {
	Id       uuid.UUID `json:"id,omitempty"`
	Name     string    `json:"name,omitempty"`
	Created  time.Time `json:"created,omitempty"`
	Modified time.Time `json:"modified,omitempty"`
}

func (l *LibraryDTO) FromModel(m model.Library) *LibraryDTO {
	l.Id = m.ID
	l.Name = m.Name
	l.Created = m.Created
	l.Modified = m.Modified

	return l
}

type LibraryUpdateDTO struct {
	Name string `json:"name"`
}
