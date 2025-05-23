package models

import "github.com/google/uuid"

type ScanPathData struct {
	LibraryPathId uuid.UUID `json:"libraryPathId"`
}

type GenerateThumbnailData struct {
	VideoId       uuid.UUID `json:"videoId"`
	LibraryPathId uuid.UUID `json:"libraryPathId"`
	Path          string    `json:"path"`
	// Optional: If set to 0, timestamp at 25% of video playback will be used
	Timestamp int `json:"timestamp"`
	// Optional: If set to 0, video height will be used
	Height int `json:"height"`
	// Optional: If set to 0, video widtch will be used
	Width int `json:"width"`
}
