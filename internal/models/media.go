package models

import (
	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
)

type MediaOrdinal string

const (
	MediaOrdinal_Created  MediaOrdinal = "created"
	MediaOrdinal_Modified MediaOrdinal = "modified"
	MediaOrdinal_Path     MediaOrdinal = "path"
	MediaOrdinal_Title    MediaOrdinal = "title"
	MediaOrdinal_Size     MediaOrdinal = "size"
	MediaOrdinal_Added    MediaOrdinal = "added"
)

func (o MediaOrdinal) ToColumn() postgres.Column {
	media := table.Media
	switch o {
	case MediaOrdinal_Created:
		return media.Created
	case MediaOrdinal_Modified:
		return media.Modified
	case MediaOrdinal_Path:
		return media.Path
	case MediaOrdinal_Title:
		return media.Title
	case MediaOrdinal_Size:
		return media.Size
	case MediaOrdinal_Added:
		return media.Added
	default:
		return media.Added
	}
}

type MediaSearchDTO struct {
	Limit   int          `form:"limit"`
	Skip    int          `form:"skip"`
	OrderBy MediaOrdinal `form:"orderBy"`
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
