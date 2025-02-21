package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
)

type CreateLibraryModel struct {
	Name string `json:"name" binding:"required"`
}

type Library struct {
	Id       uuid.UUID `json:"id,omitempty"`
	Name     string    `json:"name,omitempty"`
	Created  time.Time `json:"created,omitempty"`
	Modified time.Time `json:"modified,omitempty"`
}

func (l *Library) FromModel(m model.Library) *Library {
	l.Id = m.ID
	l.Name = m.Name
	l.Created = m.Created
	l.Modified = m.Modified

	return l
}
