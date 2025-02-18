package models

import "github.com/google/uuid"

type ScanPathData struct {
	LibraryPathId uuid.UUID `json:"libraryPathId"`
}
