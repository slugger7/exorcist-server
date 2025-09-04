package models

import (
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WSConn struct {
	Conn *websocket.Conn
	Mu   sync.Mutex
}
type WebSocketMap = map[uuid.UUID][]*WSConn
