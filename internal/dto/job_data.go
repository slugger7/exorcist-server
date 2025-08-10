package dto

import (
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
)

type ScanPathData struct {
	LibraryPathId uuid.UUID `json:"libraryPathId"`
}

type GenerateThumbnailData struct {
	MediaId uuid.UUID `json:"mediaId"`
	Path    string    `json:"path"`
	// Optional: If set to 0, timestamp at 25% of video playback will be used. Value in seconds
	Timestamp float64 `json:"timestamp"`
	// Optional: If set to 0, video height will be used
	Height int `json:"height"`
	// Optional: If set to 0, video widtch will be used
	Width        int                          `json:"width"`
	RelationType *model.MediaRelationTypeEnum `json:"relationType"`
	Metadata     *ThumbnailMetadataDTO        `json:"metadata"`
}

type RefreshFields struct {
	Size     bool `json:"size"`
	Checksum bool `json:"checksum"`
}

type RefreshMetadata struct {
	MediaId       uuid.UUID      `json:"mediaId"`
	RefreshFields *RefreshFields `json:"refreshFields"`
}

type RefreshLibraryMetadata struct {
	LibraryId     uuid.UUID      `json:"libraryId"`
	BatchSize     int            `json:"batchSize"`
	RefreshFields *RefreshFields `json:"refreshFields"`
}

type GenerateChaptersData struct {
	MediaId      uuid.UUID             `json:"mediaId"`
	Interval     float64               `json:"interval"`
	Height       int                   `json:"height"`
	Width        int                   `json:"width"`
	MaxDimension int                   `json:"maxDimension"`
	Metadata     *ChapterMetadadataDTO `json:"metadata"`
}
