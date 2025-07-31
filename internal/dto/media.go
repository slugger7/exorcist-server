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
	MediaOrdinal_Runtime  MediaOrdinal = "runtime"
)

var MediaOrdinalAllValues = []MediaOrdinal{
	MediaOrdinal_Created,
	MediaOrdinal_Modified,
	MediaOrdinal_Added,
	MediaOrdinal_Path,
	MediaOrdinal_Size,
	MediaOrdinal_Title,
}

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
	case MediaOrdinal_Runtime:
		return table.Video.Runtime
	default:
		return media.Added
	}
}

func (o MediaOrdinal) String() string {
	return string(o)
}

type WatchStatus string

const (
	WatchStatus_Watched    WatchStatus = "watched"
	WatchStatus_Unwatched  WatchStatus = "unwatched"
	WatchStatus_InProgress WatchStatus = "in_progress"
)

var WatchStatusAllValues = []WatchStatus{
	WatchStatus_Watched,
	WatchStatus_Unwatched,
	WatchStatus_InProgress,
}

func (w WatchStatus) String() string {
	return string(w)
}

type MediaSearchDTO struct {
	PageRequestDTO
	OrderBy       MediaOrdinal  `form:"orderBy" json:"orderBy"`
	Search        string        `form:"search" json:"search"`
	Tags          []string      `form:"tags" json:"tags"`
	People        []string      `form:"people" json:"people"`
	WatchStatuses []WatchStatus `form:"watchStatuses" json:"watchStatus"`
}

type MediaOverviewDTO struct {
	Id          uuid.UUID `json:"id"`
	Title       string    `json:"title,omitempty"`
	ThumbnailId uuid.UUID `json:"thumbnailId,omitempty"`
	Progress    float64   `json:"progress,omitempty"`
	Deleted     bool      `json:"deleted"`
}

func (v *MediaOverviewDTO) FromModel(m models.MediaOverviewModel) *MediaOverviewDTO {
	v.Id = m.Media.ID
	v.Title = m.Title
	v.Deleted = m.Deleted
	v.Progress = m.MediaProgress.Timestamp
	v.ThumbnailId = m.Thumbnail.ID

	return v
}

type MediaDTO struct {
	ID            uuid.UUID   `json:"id"`
	LibraryPathID uuid.UUID   `json:"libraryPathId"`
	Path          string      `json:"path"`
	Title         string      `json:"title"`
	Size          int64       `json:"size"`
	Checksum      *string     `json:"checksum"`
	Exists        bool        `json:"exists"`
	Deleted       bool        `json:"deleted"`
	Added         time.Time   `json:"added"`
	Created       time.Time   `json:"created"`
	Modified      time.Time   `json:"modified"`
	Image         *ImageDTO   `json:"image,omitempty"`
	Video         *VideoDTO   `json:"video,omitempty"`
	ThumbnailID   uuid.UUID   `json:"thumbnailId,omitempty"`
	Progress      float64     `json:"progress"`
	People        []PersonDTO `json:"people"`
	Tags          []TagDTO    `json:"tags"`
}

func (d *MediaDTO) FromModel(m models.Media) *MediaDTO {
	d.ID = m.Media.ID
	d.LibraryPathID = m.LibraryPathID
	d.Path = m.Path
	d.Title = m.Title
	d.Size = m.Size
	d.Checksum = m.Checksum
	d.Deleted = m.Deleted
	d.Exists = m.Exists
	d.Added = m.Added
	d.Created = m.Media.Created
	d.Modified = m.Media.Modified

	if m.MediaProgress != nil {
		d.Progress = m.MediaProgress.Timestamp
	}

	if m.Thumbnail != nil {
		d.ThumbnailID = m.Thumbnail.ID
	}

	d.Image = (&ImageDTO{}).FromModel(m.Image)
	d.Video = (&VideoDTO{}).FromModel(m.Video)

	if len(m.People) > 0 {
		d.People = make([]PersonDTO, len(m.People))
		for i, p := range m.People {
			d.People[i] = *(&PersonDTO{}).FromModel(&p)
		}
	}

	if len(m.Tags) > 0 {
		d.Tags = make([]TagDTO, len(m.Tags))
		for i, t := range m.Tags {
			d.Tags[i] = *(&TagDTO{}).FromModel(&t)
		}
	}

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

type DeleteMediaDTO struct {
	Physical *bool `json:"physical" form:"physical"`
}

type MediaUpdateDTO struct {
	Title *string `json:"title" form:"title"`
}

type MediaUpdatedDTO struct {
	ID       uuid.UUID `json:"id" form:"id"`
	Title    *string   `json:"title,omitempty" form:"title"`
	Modified time.Time `json:"modified" form:"modified"`
}

func (d *MediaUpdatedDTO) FromModel(m model.Media) *MediaUpdatedDTO {
	d.ID = m.ID
	d.Title = &m.Title
	d.Modified = m.Modified

	return d
}
