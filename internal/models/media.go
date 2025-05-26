package models

import (
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
)

type MediaSearchDTO struct {
	Limit   int          `form:"limit"`
	Skip    int          `form:"skip"`
	OrderBy VideoOrdinal `form:"orderBy"`
	Asc     bool         `form:"asc"`
	Search  string       `form:"search"`
}

type MediaVideo struct {
	model.Media
	model.Video
}

type MediaOverviewModel struct {
	model.Media
	ThumbnailId uuid.UUID
}

type MediaOverviewDTO struct {
	Id          uuid.UUID `json:"id"`
	Title       string    `json:"title,omitempty"`
	ThumbnailId uuid.UUID `json:"thumbnailId,omitempty"`
	Deleted     bool      `json:"deleted"`
}

func (v *MediaOverviewModel) ToDTO() *MediaOverviewDTO {
	return &MediaOverviewDTO{
		Id:          v.Media.ID,
		Title:       v.Media.Title,
		ThumbnailId: v.ThumbnailId,
		Deleted:     v.Deleted,
	}
}

func (v *MediaOverviewDTO) FromModel(m *model.Media, i *model.Media) *MediaOverviewDTO {
	v.Id = m.ID
	v.Title = m.Title
	v.Deleted = m.Deleted
	if i != nil {
		v.ThumbnailId = i.ID
	}
	return v
}
