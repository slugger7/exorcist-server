package dto

import "github.com/google/uuid"

type ScanPathData struct {
	LibraryPathId uuid.UUID `json:"libraryPathId"`
}

type GenerateThumbnailData struct {
	VideoId uuid.UUID `json:"videoId"`
	Path    string    `json:"path"`
	// Optional: If set to 0, timestamp at 25% of video playback will be used
	Timestamp int `json:"timestamp"`
	// Optional: If set to 0, video height will be used
	Height int `json:"height"`
	// Optional: If set to 0, video widtch will be used
	Width int `json:"width"`
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
