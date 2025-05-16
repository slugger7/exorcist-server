package models

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	errs "github.com/slugger7/exorcist/internal/errors"
)

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

type WebSocketMap = map[uuid.UUID][]*websocket.Conn

func (msg *WSMessage[T]) SendToAll(wss WebSocketMap) error {
	for _, ws := range wss {
		for _, s := range ws {
			if err := s.WriteJSON(msg); err != nil {
				return errs.BuildError(err, "could not write json to websocket")
			}
		}
	}

	return nil
}
