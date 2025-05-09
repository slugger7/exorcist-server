package models

import (
	"github.com/go-jet/jet/v2/postgres"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
)

type VideoOrdinal string

const (
	VideoOrdinal_Created  VideoOrdinal = "created"
	VideoOrdinal_Modified              = "modified"
	VideoOrdinal_Title                 = "title"
	VideoOrdinal_Size                  = "size"
	VideoOrdinal_Runtime               = "runtime"
	VideoOrdinal_Added                 = "added"
)

func (o VideoOrdinal) ToColumn() postgres.Column {
	switch o {
	case VideoOrdinal_Created:
		return table.Video.Created
	case VideoOrdinal_Modified:
		return table.Video.Modified
	case VideoOrdinal_Title:
		return table.Video.Title
	case VideoOrdinal_Size:
		return table.Video.Size
	case VideoOrdinal_Runtime:
		return table.Video.Runtime
	case VideoOrdinal_Added:
		return table.Video.Added
	default:
		return table.Video.Added
	}
}
