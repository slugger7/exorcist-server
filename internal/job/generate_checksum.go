package job

import (
	"database/sql"
	"encoding/json"
	"log"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/media"
	videoRepository "github.com/slugger7/exorcist/internal/repository/video"
)

func GenerateChecksums(db *sql.DB) {
	selectStatement := videoRepository.GetVideoWithoutChecksumStatement()
	var data []struct {
		model.LibraryPath
		model.Video
	}
	err := selectStatement.Query(db, &data)
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
			Exec(db)
		if err != nil {
			log.Printf("Could not update the checksum of video (%v): %v", v.Video.ID, err)
		}
	}
}

type GenerateChecksumData struct {
	VideoId uuid.UUID `json:"videoId"`
}

func GenerateChecksum(db *sql.DB, job model.Job) {
	var jobData *GenerateChecksumData
	err := json.Unmarshal([]byte(*job.Data), &jobData)
	if err != nil {
		log.Printf("Could not read json data %s", err)
		return
	}
}
