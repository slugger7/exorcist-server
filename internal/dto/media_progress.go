package dto

import "github.com/slugger7/exorcist/internal/db/exorcist/public/model"

type ProgressDTO struct {
	Progress float64 `json:"progress"`
}

func (d *ProgressDTO) FromModel(m model.MediaProgress) *ProgressDTO {
	d.Progress = m.Timestamp
	return d
}

type ProgressUpdateDTO struct {
	Overwrite bool    `json:"overwrite" form:"overwrite"`
	Progress  float64 `json:"progress" form:"progress"`
}
