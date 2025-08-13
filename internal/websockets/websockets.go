package websockets

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/dto"
	"github.com/slugger7/exorcist/internal/environment"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/models"
)

type Websockets interface {
	AddWs(id uuid.UUID, wsConn *models.WSConn)
	WebSocketHeartbeat()
	PingDuration() time.Duration
	PongDuration() time.Duration
	Shutdown()

	MediaOverviewUpdate(media dto.MediaOverviewDTO)
	MediaUpdate(media dto.MediaDTO)
	MediaDelete(media dto.MediaOverviewDTO)
	MediaCreate(media dto.MediaOverviewDTO)
	JobUpdate(job model.Job)
	JobCreate(job model.Job)
}

type websockets struct {
	env     *environment.EnvironmentVariables
	wss     models.WebSocketMap
	logger  logger.Logger
	wsMutex sync.Mutex
}

// Shutdown implements Websockets.
func (w *websockets) Shutdown() {
	w.wsMutex.Lock()
	defer w.wsMutex.Unlock()

	w.logger.Debug("Closing websockets")
	for _, i := range w.wss {
		for _, s := range i {
			s.Mu.Lock()
			defer s.Mu.Unlock()
			s.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			s.Conn.Close()
		}
	}

	w.logger.Debug("Websockets closed")
}

// AddWs implements Websockets.
func (w *websockets) AddWs(id uuid.UUID, wsConn *models.WSConn) {
	w.wsMutex.Lock()
	w.wss[id] = append(w.wss[id], wsConn)
	w.wsMutex.Unlock()

}

func (s *websockets) PongDuration() time.Duration {
	return time.Duration(s.env.WebsocketHeartbeatInterval * int(time.Millisecond))
}

func (s *websockets) PingDuration() time.Duration {
	val := (s.env.WebsocketHeartbeatInterval * 9) / 10
	return time.Duration(val * int(time.Millisecond))
}

func (s *websockets) WebSocketHeartbeat() {
	tickerDuration := s.PingDuration()
	ticker := time.NewTicker(tickerDuration)
	defer ticker.Stop()

	for range ticker.C {
		for i, ws := range s.wss {
			for _, c := range ws {
				c.Mu.Lock()

				c.Conn.SetWriteDeadline(time.Now().Add(s.PongDuration()))
				s.logger.Debug("protocol ping")
				if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					s.logger.Warningf("could not write ping message to %v: %v", i, err.Error())
					s.wsMutex.Lock()
					delete(s.wss, i)
					s.wsMutex.Unlock()
				}
				c.Mu.Unlock()
			}
		}
	}
}

var websocketsInterface *websockets

func New(env *environment.EnvironmentVariables) Websockets {
	if websocketsInterface == nil {
		websocketsInterface = &websockets{
			env:    env,
			wss:    make(models.WebSocketMap),
			logger: logger.New(env),
		}
	}

	return websocketsInterface
}
