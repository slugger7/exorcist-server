package models

type WSTopic = string

const (
	WSTopic_JobUpdate   WSTopic = "job_update"
	WSTopic_JobCreate   WSTopic = "job_create"
	WSTopic_VideoUpdate WSTopic = "video_update"
	WSTopic_VideoCreate WSTopic = "video_create"
	WSTopic_VideoDelete WSTopic = "video_delete"
)

type WSMessage[T any] struct {
	Topic WSTopic `json:"topic"`
	Data  T       `json:"data,omitempty"`
}
