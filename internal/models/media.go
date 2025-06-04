package models

import (
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
)

type MediaVideo struct {
	model.Media
	model.Video
}

type MediaOverviewModel struct {
	model.Media
	Thumbnail
}

type Thumbnail struct {
	ID uuid.UUID `sql:"primary_key" json:"id"`
}

type Media struct {
	model.Media
	*model.Image
	*model.Video
	*Thumbnail
	People []model.Person
}
