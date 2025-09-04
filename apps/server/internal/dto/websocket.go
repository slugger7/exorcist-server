package dto

import (
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/models"
)

type WSTopic string

const (
	WSTopic_JobUpdate           WSTopic = "job_update"
	WSTopic_JobCreate           WSTopic = "job_create"
	WSTopic_MediaUpdate         WSTopic = "media_update"
	WSTopic_MediaOverviewUpdate WSTopic = "media_overview_update"
	WSTopic_MediaCreate         WSTopic = "media_create"
	WSTopic_MediaDelete         WSTopic = "media_delete"
)

var WSTopicAllValues = []WSTopic{
	WSTopic_JobUpdate,
	WSTopic_JobCreate,
	WSTopic_MediaUpdate,
	WSTopic_MediaOverviewUpdate,
	WSTopic_MediaCreate,
	WSTopic_MediaDelete,
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
