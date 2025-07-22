package server

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/slugger7/exorcist/internal/models"
)

func (s *server) pongDuration() time.Duration {
	return time.Duration(s.env.WebsocketHeartbeatInterval * int(time.Millisecond))
}

func (s *server) pingDuration() time.Duration {
	val := (s.env.WebsocketHeartbeatInterval * 9) / 10
	return time.Duration(val * int(time.Millisecond))
}

func (s *server) webSocketHeartbeat() {
	tickerDuration := s.pingDuration()
	ticker := time.NewTicker(tickerDuration)
	defer ticker.Stop()

	for range ticker.C {
		for i, ws := range s.websockets {
			for _, c := range ws {
				c.Mu.Lock()

				c.Conn.SetWriteDeadline(time.Now().Add(s.pongDuration()))
				s.logger.Debug("protocol ping")
				if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					s.logger.Warningf("could not write ping message to %v: %v", i, err.Error())
					s.websocketMutex.Lock()
					delete(s.websockets, i)
					s.websocketMutex.Unlock()
				}
				c.Mu.Unlock()
			}
		}
	}
}

func (s *server) withWS(r *gin.RouterGroup, route Route) *server {
	r.GET(fmt.Sprintf("%v/ws", route), s.ws)

	go s.webSocketHeartbeat()
	return s
}

func (s *server) ws(c *gin.Context) {
	session := sessions.Default(c)
	if userRaw, ok := session.Get(userKey).(string); ok {
		userId, err := uuid.Parse(userRaw)
		if err != nil {
			s.logger.Warningf("could not parse user id from string: %v", userRaw)
			c.AbortWithStatus(http.StatusUnprocessableEntity)
		}

		upgrader := s.wsUpgrader()
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			s.logger.Infof("Origin: %v", c.Request.Host)
			s.logger.Errorf("failed to upgrade connection to web socket: %v", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		conn.SetPongHandler(func(string) error {
			conn.SetReadDeadline(time.Now().Add(s.pongDuration()))
			s.logger.Debug("protocol pong")
			return nil
		})

		wsConn := models.WSConn{
			Conn: conn, Mu: sync.Mutex{},
		}

		s.websocketMutex.Lock()
		s.websockets[userId] = append(s.websockets[userId], &wsConn)
		s.websocketMutex.Unlock()

		go s.wsReader(conn, userId)
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

func (s *server) wsReader(ws *websocket.Conn, id uuid.UUID) {
	for {
		ws.SetReadDeadline(time.Now().Add(s.pongDuration()))
		_, message, err := ws.ReadMessage()
		if err != nil {
			s.logger.Errorf("could not read message: %v", err.Error())
			break
		}

		message = bytes.TrimSpace(bytes.Replace(message, []byte{'\n'}, []byte{' '}, -1))
		if string(message) == "ping" {
			s.logger.Debug("application pong")
			ws.WriteMessage(websocket.TextMessage, []byte("pong"))
		}
	}
}
