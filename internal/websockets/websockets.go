package websockets

import (
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/dto"
	"github.com/slugger7/exorcist/internal/environment"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/models"
)

type Websockets interface {
	MediaOverviewUpdate(media dto.MediaOverviewDTO)
	MediaUpdate(media dto.MediaDTO)
	MediaDelete(media dto.MediaOverviewDTO)
	MediaCreate(media dto.MediaOverviewDTO)
	JobUpdate(job model.Job)
}

type websockets struct {
	env    *environment.EnvironmentVariables
	wss    models.WebSocketMap
	logger logger.Logger
}

var websocketsInterface *websockets

func New(env *environment.EnvironmentVariables, wss models.WebSocketMap) Websockets {
	if websocketsInterface == nil {
		websocketsInterface = &websockets{
			env:    env,
			wss:    wss,
			logger: logger.New(env),
		}
	}

	return websocketsInterface
}
