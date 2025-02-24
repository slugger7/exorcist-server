package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
)

type CreateLibraryPathModel struct {
	LibraryId uuid.UUID `json:"libraryId" binding:"required"`
	Path      string    `json:"path" binding:"required"`
}

type LibraryPathModel struct {
	Id        uuid.UUID `json:"id,omitempty"`
	LibraryId uuid.UUID `json:"libraryId,omitempty"`
	Path      string    `json:"path,omitempty"`
	Created   time.Time `json:"created,omitempty"`
	Modified  time.Time `json:"modified,omitempty"`
}

func (l *LibraryPathModel) FromModel(m model.LibraryPath) *LibraryPathModel {
	l.Id = m.ID
	l.LibraryId = m.LibraryID
	l.Path = m.Path
	l.Created = m.Created
	l.Modified = m.Modified

	return l
}
