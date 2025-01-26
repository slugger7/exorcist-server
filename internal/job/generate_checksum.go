package job

import (
	"database/sql"
	"log"
	"path/filepath"

	"github.com/slugger7/exorcist/internal/media"
	videoRepository "github.com/slugger7/exorcist/internal/repository/video"
)

func GenerateChecksums(db *sql.DB) {
	selectStatement := videoRepository.GetVideoWithoutChecksumStatement()
	data, err := videoRepository.ExecuteChecksumStatement(db, selectStatement)
	if err != nil {
		log.Printf("Error while fetching a video without a checksum %v", err.Error())
	}

	for _, v := range data {
		absolutePath := filepath.Join(v.LibraryPath.Path, v.Video.RelativePath)
		checksum, err := media.CalculateMD5(absolutePath)
		if err != nil {
			log.Printf("Could not calculate checksum for %v. Video ID %v", absolutePath, v.Video.ID)
		}

		v.Video.Checksum = &checksum

		videoRepository.ExecuteUpdate(db, videoRepository.UpdateVideoChecksum(v.Video))
	}
}
