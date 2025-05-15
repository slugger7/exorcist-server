package models

import (
	"strconv"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
)

type VideoOverviewModel struct {
	model.Video
	model.Image
	model.LibraryPath
}

type VideoOverviewDTO struct {
	Id    uuid.UUID `json:"id"`
	Title string    `json:"title,omitempty"`
	// Deprecated: this is useless on the frontend at the moment
	Path        string    `json:"path,omitempty"`
	ThumbnailId uuid.UUID `json:"thumbnailId,omitempty"`
}

func (v *VideoOverviewModel) ToDTO() *VideoOverviewDTO {
	return &VideoOverviewDTO{
		Id:          v.Video.ID,
		Title:       v.Video.Title,
		Path:        v.LibraryPath.Path + v.Video.RelativePath,
		ThumbnailId: v.Image.ID,
	}
}

func (v *VideoOverviewDTO) FromModel(m *model.Video, i *model.Image) *VideoOverviewDTO {
	v.Id = m.ID
	v.Title = m.Title
	if i != nil {
		v.ThumbnailId = i.ID
	}
	return v
}

func DefualtBool(strVal string, def bool) bool {
	val, err := strconv.ParseBool(strVal)
	if err != nil {
		return def
	}
	return val
}

func DefaultInt(strVal string, def int) int {
	if strVal != "" {
		val, err := strconv.Atoi(strVal)
		if err == nil {
			return val
		}
	}

	return def
}

type VideoSearchDTO struct {
	Limit   int          `form:"limit"`
	Skip    int          `form:"skip"`
	OrderBy VideoOrdinal `form:"orderBy"`
	Asc     bool         `form:"asc"`
	Search  string       `form:"search"`
}
