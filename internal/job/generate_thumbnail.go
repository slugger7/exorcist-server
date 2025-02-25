package job

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/ffmpeg"
)

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

func (jr *JobRunner) GenerateThumbnail(job *model.Job) error {
	var jobData GenerateThumbnailData
	if err := json.Unmarshal([]byte(*job.Data), &jobData); err != nil {
		return errs.BuildError(err, "error parsing job data: %v", job.Data)
	}

	if jobData.Path == "" {
		return fmt.Errorf("cant create an image at a blank path")
	}

	video, err := jr.repo.Video().GetByIdWithLibraryPath(jobData.VideoId)
	if err != nil {
		return errs.BuildError(err, "error fetching video with library path by id: %v", jobData.VideoId)
	}

	if jobData.Height == 0 {
		jobData.Height = int(video.Height)
	}
	if jobData.Width == 0 {
		jobData.Width = int(video.Width)
	}
	if jobData.Timestamp == 0 {
		jobData.Timestamp = int(float64(video.Runtime) * 0.25)
	}

	absolutePath := filepath.Join(video.LibraryPath.Path, video.RelativePath)

	if err := ffmpeg.ImageAt(absolutePath, jobData.Timestamp, jobData.Path, jobData.Width, jobData.Height); err != nil {
		return errs.BuildError(err, "could not create image at timestamp")
	}

	image := &model.Image{
		Name: video.Title,
		Path: jobData.Path,
	}

	image, err = jr.repo.Image().Create(image)
	if err != nil {
		return errs.BuildError(err, "error creating image")
	}

	videoImage := &model.VideoImage{
		VideoID:        video.Video.ID,
		ImageID:        image.ID,
		VideoImageType: model.VideoImageTypeEnum_Thumbnail,
	}

	videoImage, err = jr.repo.Image().RelateVideo(videoImage)
	if err != nil {
		return errs.BuildError(err, "could not create video image relation")
	}

	return nil
}
