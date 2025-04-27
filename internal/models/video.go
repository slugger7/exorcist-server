package models

import (
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
)

type VideoOverviewModel struct {
	model.Video
	model.Image
	model.LibraryPath
}

type VideoOverviewDTO struct {
	Id            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	Path          string    `json:"path"`
	ThumbnailPath string    `json:"thumbnailPath"`
}

func (v *VideoOverviewModel) ToDTO() VideoOverviewDTO {
	return VideoOverviewDTO{
		Id:            v.Video.ID,
		Title:         v.Video.Title,
		Path:          v.LibraryPath.Path + v.Video.RelativePath,
		ThumbnailPath: v.Image.Path,
	}
}
