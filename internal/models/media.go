package models

import "github.com/slugger7/exorcist/internal/db/exorcist/public/model"

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
