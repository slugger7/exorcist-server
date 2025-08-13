package websockets

import (
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/dto"
)

// JobUpdate implements Websockets.
func (w *websockets) JobUpdate(job model.Job) {
	w.logger.Debug("ws - updating job")

	jobUpdate := dto.WSMessage[dto.JobDTO]{
		Topic: dto.WSTopic_JobUpdate,
		Data:  *(&dto.JobDTO{}).FromModel(job),
	}
	jobUpdate.SendToAll(w.wss)
}

// JobCreate implements Websockets.
func (w *websockets) JobCreate(job model.Job) {
	w.logger.Debug("ws - createing job")

	jobDto := (&dto.JobDTO{}).FromModel(job)
	message := dto.WSMessage[dto.JobDTO]{
		Topic: dto.WSTopic_JobCreate,
		Data:  *jobDto,
	}

	message.SendToAll(w.wss)
}
