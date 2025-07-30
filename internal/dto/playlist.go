package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
)

type CreatePlaylistDTO struct {
	Name string `json:"string"`
}

type PlaylistDTO struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}

func (d *PlaylistDTO) FromModel(m model.Playlist) *PlaylistDTO {
	d.ID = m.ID
	d.Name = m.Name
	d.Created = m.Created
	d.Modified = m.Modified

	return d
}

type CreatePlaylistMediaDTO struct {
	MediaID uuid.UUID `json:"mediaId"`
}
