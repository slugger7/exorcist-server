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
	model.MediaProgress
	*model.Video
	Thumbnail
	*model.FavouriteMedia
}

type Thumbnail struct {
	ID uuid.UUID `sql:"primary_key" json:"id"`
}

type Media struct {
	model.Media
	*model.Image
	*model.Video
	*Thumbnail
	*model.MediaProgress
	*model.FavouriteMedia
	People []model.Person
	Tags   []model.Tag
}
