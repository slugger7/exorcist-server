package dto

import (
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/models"
)

type MediaOrdinal string

const (
	MediaOrdinal_Created  MediaOrdinal = "created"
	MediaOrdinal_Modified MediaOrdinal = "modified"
	MediaOrdinal_Path     MediaOrdinal = "path"
	MediaOrdinal_Title    MediaOrdinal = "title"
	MediaOrdinal_Size     MediaOrdinal = "size"
	MediaOrdinal_Added    MediaOrdinal = "added"
)

func (o MediaOrdinal) ToColumn() postgres.Column {
	media := table.Media
	switch o {
	case MediaOrdinal_Created:
		return media.Created
	case MediaOrdinal_Modified:
		return media.Modified
	case MediaOrdinal_Path:
		return media.Path
	case MediaOrdinal_Title:
		return media.Title
	case MediaOrdinal_Size:
		return media.Size
	case MediaOrdinal_Added:
		return media.Added
	default:
		return media.Added
	}
}

type MediaSearchDTO struct {
	PageRequestDTO
	OrderBy MediaOrdinal `form:"orderBy" json:"orderBy"`
	Search  string       `form:"search" json:"search"`
}

type MediaOverviewDTO struct {
	Id          uuid.UUID `json:"id"`
	Title       string    `json:"title,omitempty"`
	ThumbnailId uuid.UUID `json:"thumbnailId,omitempty"`
	Deleted     bool      `json:"deleted"`
}

func (v *MediaOverviewDTO) FromModel(m *model.Media, i *model.Media) *MediaOverviewDTO {
	v.Id = m.ID
	v.Title = m.Title
	v.Deleted = m.Deleted
	if i != nil {
		v.ThumbnailId = i.ID
	}
	return v
}

type MediaDTO struct {
	ID            uuid.UUID `json:"id"`
	LibraryPathID uuid.UUID `json:"libraryPathId"`
	Path          string    `json:"path"`
	Title         string    `json:"title"`
	Size          int64     `json:"size"`
	Checksum      *string   `json:"checksum"`
	Added         time.Time `json:"added"`
	Created       time.Time `json:"created"`
	Modified      time.Time `json:"modified"`
	Image         *ImageDTO `json:"image,omitempty"`
	Video         *VideoDTO `json:"video,omitempty"`
	ThumbnailID   uuid.UUID `json:"thumbnailId,omitempty"`
}

func (d *MediaDTO) FromModel(m models.Media) *MediaDTO {
	d.ID = m.Media.ID
	d.LibraryPathID = m.LibraryPathID
	d.Path = m.Path
	d.Title = m.Title
	d.Size = m.Size
	d.Checksum = m.Checksum
	d.Added = m.Added
	d.Created = m.Created
	d.Modified = m.Modified
	d.ThumbnailID = m.Thumbnail.ID

	d.Image = (&ImageDTO{}).FromModel(m.Image)
	d.Video = (&VideoDTO{}).FromModel(m.Video)

	return d
}

type ImageDTO struct {
	ID      uuid.UUID `json:"id"`
	MediaID uuid.UUID `json:"mediaId"`
	Height  int32     `json:"height"`
	Width   int32     `json:"width"`
}

func (d *ImageDTO) FromModel(m *model.Image) *ImageDTO {
	if m == nil {
		return nil
	}

	d.ID = m.ID
	d.MediaID = m.MediaID
	d.Height = m.Height
	d.Width = m.Width

	return d
}

type VideoDTO struct {
	ID      uuid.UUID `json:"id"`
	MediaID uuid.UUID `json:"mediaId"`
	Height  int32     `json:"height"`
	Width   int32     `json:"width"`
	Runtime float64   `json:"runtime"`
}

func (d *VideoDTO) FromModel(m *model.Video) *VideoDTO {
	if m == nil {
		return nil
	}

	d.ID = m.ID
	d.MediaID = m.MediaID
	d.Height = m.Height
	d.Width = m.Width
	d.Runtime = m.Runtime

	return d
}
