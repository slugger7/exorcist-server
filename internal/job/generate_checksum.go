package job

import (
	"log"
	"path/filepath"

	"github.com/slugger7/exorcist/internal/db"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/media"
	videoRepository "github.com/slugger7/exorcist/internal/repository/video"
)

func GenerateChecksums(db db.Service) {
	selectStatement := videoRepository.GetVideoWithoutChecksumStatement()
	var data []struct {
		model.LibraryPath
		model.Video
	}
	err := selectStatement.Query(db.GetDb(), &data)
	if err != nil {
		log.Printf("Error while fetching a video without a checksum %v", err.Error())
	}

	for _, v := range data {
		absolutePath := filepath.Join(v.LibraryPath.Path, v.Video.RelativePath)
		log.Printf("Calculating checksum for %v", v.Video.RelativePath)
		checksum, err := media.CalculateMD5(absolutePath)
		if err != nil {
			log.Printf("Could not calculate checksum for %v. Video ID %v", absolutePath, v.Video.ID)
		}

		v.Video.Checksum = &checksum

		_, err = videoRepository.UpdateVideoChecksum(v.Video).
			Exec(db.GetDb())
		if err != nil {
			log.Printf("Could not update the checksum of video (%v): %v", v.Video.ID, err)
		}
	}
	log.Println("Completed checksum generation")
}
