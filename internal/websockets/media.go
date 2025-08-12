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
	logger logger.ILogger
}

// JobUpdate implements Websockets.
func (w *websockets) JobUpdate(job model.Job) {
	w.logger.Debug("ws - updating job")

	jobUpdate := dto.WSMessage[dto.JobDTO]{
		Topic: dto.WSTopic_JobUpdate,
		Data:  *(&dto.JobDTO{}).FromModel(job),
	}
	jobUpdate.SendToAll(w.wss)
}

// MediaCreate implements Websockets.
func (w *websockets) MediaCreate(media dto.MediaOverviewDTO) {
	w.logger.Debug("ws - creating video")

	mediaDelete := dto.WSMessage[dto.MediaOverviewDTO]{
		Topic: dto.WSTopic_MediaCreate,
		Data:  media,
	}
	mediaDelete.SendToAll(w.wss)

}

// MediaDelete implements Websockets.
func (w *websockets) MediaDelete(media dto.MediaOverviewDTO) {
	w.logger.Debug("ws - deleting video")

	videoDelete := dto.WSMessage[dto.MediaOverviewDTO]{
		Topic: dto.WSTopic_MediaDelete,
		Data:  media,
	}
	videoDelete.SendToAll(w.wss)
}

// MediaUpdate implements Websockets.
func (w *websockets) MediaUpdate(media dto.MediaDTO) {
	w.logger.Debug("ws - updating media")

	mediaUpdate := dto.WSMessage[dto.MediaDTO]{
		Topic: dto.WSTopic_MediaUpdate,
		Data:  media,
	}

	mediaUpdate.SendToAll(w.wss)
}

// MediaOverviewUpdate implements Websockets.
func (w *websockets) MediaOverviewUpdate(media dto.MediaOverviewDTO) {
	w.logger.Debug("ws - updating media overview")

	mediaUpdate := dto.WSMessage[dto.MediaOverviewDTO]{
		Topic: dto.WSTopic_MediaOverviewUpdate,
		Data:  media,
	}
	mediaUpdate.SendToAll(w.wss)
}

var websocketsInterface *websockets

func New(env *environment.EnvironmentVariables, wss models.WebSocketMap) Websockets {
	if websocketsInterface == nil {
		websocketsInterface = &websockets{
			env: env,
			wss: wss,
		}
	}

	return websocketsInterface
}
