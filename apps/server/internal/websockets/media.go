package websockets

import (
	"github.com/slugger7/exorcist/internal/dto"
)

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
