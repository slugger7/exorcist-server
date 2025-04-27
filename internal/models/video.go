package models

import "github.com/google/uuid"

type VideoOverviewDTO struct {
	Id            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	Path          string    `json:"path"`
	ThumbnailPath string    `json:"thumbnailPath"`
}
