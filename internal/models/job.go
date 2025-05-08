package models

import "github.com/slugger7/exorcist/internal/db/exorcist/public/model"

type CreateJobDTO struct {
	Type     model.JobTypeEnum      `json:"type" binding:"required"`
	Data     map[string]interface{} `json:"data"`
	Priority *int16                 `json:"priority"`
}

type JobPriority = int16

const (
	JobPriority_Highest JobPriority = iota
	JobPriority_High
	JobPriority_MediumHigh
	JobPriority_Medium
	JobPriority_MediumLow
	JobPriority_Low
	JobPriority_Lowest
)
