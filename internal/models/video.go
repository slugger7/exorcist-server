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
	Id          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Path        string    `json:"path"`
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
		if err != nil {
			return def
		}

		return val
	}

	return def
}

type VideoSearch struct {
	Limit   int
	Skip    int
	OrderBy VideoOrdinal
	Asc     bool
	Search  string
}
