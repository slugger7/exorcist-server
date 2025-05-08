package models

import "github.com/slugger7/exorcist/internal/db/exorcist/public/model"

type CreateJobDTO struct {
	Type     model.JobTypeEnum      `json:"type" binding:"required"`
	Data     map[string]interface{} `json:"data"`
	Priority *int16                 `json:"priority"`
}
