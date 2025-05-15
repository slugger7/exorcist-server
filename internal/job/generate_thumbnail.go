package job

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/ffmpeg"
	"github.com/slugger7/exorcist/internal/models"
)

func CreateGenerateThumbnailJob(videoId, jobId uuid.UUID, imagePath string, timestamp, height, width int) (*model.Job, error) {
	d := models.GenerateThumbnailData{
		VideoId:   videoId,
		Path:      imagePath,
		Height:    height,
		Width:     width,
		Timestamp: timestamp,
	}

	js, err := json.Marshal(d)
	if err != nil {
		return nil, errs.BuildError(err, "could not marshal generate thumbnail data")
	}
	data := string(js)
	job := &model.Job{
		JobType:  model.JobTypeEnum_GenerateThumbnail,
		Status:   model.JobStatusEnum_NotStarted,
		Data:     &data,
		Parent:   &jobId,
		Priority: models.JobPriority_MediumHigh,
	}

	return job, nil
}

func createAssetDirectory(path string) error {
	dir := filepath.Dir(path)
	return os.MkdirAll(dir, os.ModePerm)
}

func (jr *JobRunner) GenerateThumbnail(job *model.Job) error {
	var jobData models.GenerateThumbnailData
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
	err = createAssetDirectory(jobData.Path)
	if err != nil {
		return errs.BuildError(err, "could not create path for asset")
	}

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

	vidUpdate := models.VideoOverviewDTO{
		Id:          video.Video.ID,
		ThumbnailId: videoImage.ImageID,
	}
	jr.wsVideoUpdate(vidUpdate)

	return nil
}
