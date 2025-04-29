package videoRepository

import (
	"github.com/go-jet/jet/v2/postgres"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
)

type VideoOrdinal string

const (
	Created  VideoOrdinal = "created"
	Modified VideoOrdinal = "modified"
	Title    VideoOrdinal = "title"
	Size     VideoOrdinal = "size"
	Runtime  VideoOrdinal = "runtime"
	Added    VideoOrdinal = "added"
)

func ordinalToColumn(ordinal *VideoOrdinal) postgres.Column {
	switch *ordinal {
	case Created:
		return table.Video.Created
	case Modified:
		return table.Video.Modified
	case Title:
		return table.Video.Title
	case Size:
		return table.Video.Size
	case Runtime:
		return table.Video.Runtime
	case Added:
		return table.Video.Added
	default:
		return table.Video.Added
	}
}
