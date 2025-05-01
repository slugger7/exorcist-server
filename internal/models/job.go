package models

import "github.com/slugger7/exorcist/internal/db/exorcist/public/model"

type CreateJob struct {
	Type model.JobTypeEnum `json:"type" binding:"required"`
	Data *string           `json:"data"`
}
