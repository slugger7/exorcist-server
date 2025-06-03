package dto

import (
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/models"
)

type WSTopic string

const (
	WSTopic_JobUpdate   WSTopic = "job_update"
	WSTopic_JobCreate   WSTopic = "job_create"
	WSTopic_VideoUpdate WSTopic = "video_update"
	WSTopic_VideoCreate WSTopic = "video_create"
	WSTopic_VideoDelete WSTopic = "video_delete"
)

var WSTopicAllValues = []WSTopic{
	WSTopic_JobUpdate,
	WSTopic_JobCreate,
	WSTopic_VideoUpdate,
	WSTopic_VideoCreate,
	WSTopic_VideoDelete,
}

func (t WSTopic) String() string {
	return string(t)
}

type WSMessage[T any] struct {
	Topic WSTopic `json:"topic"`
	Data  T       `json:"data,omitempty"`
}

func (msg *WSMessage[T]) SendToAll(wss models.WebSocketMap) error {
	for _, ws := range wss {
		for _, s := range ws {
			s.Mu.Lock()
			if err := s.Conn.WriteJSON(msg); err != nil {
				s.Mu.Unlock()
				return errs.BuildError(err, "could not write json to websocket")
			}
			s.Mu.Unlock()
		}
	}

	return nil
}
