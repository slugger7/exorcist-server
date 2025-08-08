package job

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/dto"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/ffmpeg"
	"github.com/slugger7/exorcist/internal/media"
)

func CreateGenerateThumbnailJob(video model.Video, jobId *uuid.UUID, imagePath string, timestamp, height, width int) (*model.Job, error) {
	d := dto.GenerateThumbnailData{
		MediaId:   video.ID,
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
		Parent:   jobId,
		Priority: dto.JobPriority_MediumHigh,
	}

	return job, nil
}

func createAssetDirectory(path string) error {
	dir := filepath.Dir(path)
	return os.MkdirAll(dir, os.ModePerm)
}

func (jr *JobRunner) GenerateThumbnail(job *model.Job) error {
	var jobData dto.GenerateThumbnailData
	if err := json.Unmarshal([]byte(*job.Data), &jobData); err != nil {
		return errs.BuildError(err, "error parsing job data: %v", job.Data)
	}

	if jobData.Path == "" {
		return fmt.Errorf("cant create an image at a blank path")
	}

	video, err := jr.repo.Video().GetByIdWithMedia(jobData.MediaId)
	if err != nil {
		return errs.BuildError(err, "error fetching video with id: %v", jobData.MediaId)
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

	err = createAssetDirectory(jobData.Path)
	if err != nil {
		return errs.BuildError(err, "could not create path for asset")
	}

	if err := ffmpeg.ImageAt(video.Path, jobData.Timestamp, jobData.Path, jobData.Width, jobData.Height); err != nil {
		return errs.BuildError(err, "could not create image at timestamp: %v, video: %v", jobData.Timestamp, video.Runtime)
	}

	fileSize, err := media.GetFileSize(jobData.Path)
	if err != nil {
		return errs.BuildError(err, "could not get file size for: %v", jobData.Path)
	}

	imageMedia := &model.Media{
		LibraryPathID: video.LibraryPathID,
		Path:          jobData.Path,
		Title:         fmt.Sprintf("%v-thumbnail", video.MediaID),
		MediaType:     model.MediaTypeEnum_Asset,
		Size:          fileSize,
	}

	newModels, err := jr.repo.Media().Create([]model.Media{*imageMedia})
	if err != nil {
		return errs.BuildError(err, "could not create image media")
	}
	if len(newModels) != 1 {
		return fmt.Errorf("length of models was not 1 but %v", len(newModels))
	}

	image := &model.Image{
		MediaID: newModels[0].ID,
		Height:  int32(jobData.Height),
		Width:   int32(jobData.Width),
	}

	image, err = jr.repo.Image().Create(image)
	if err != nil {
		return errs.BuildError(err, "error creating image")
	}

	videoImage := &model.MediaRelation{
		MediaID:      video.Media.ID,
		RelatedTo:    image.MediaID,
		RelationType: model.MediaRelationTypeEnum_Thumbnail,
	}

	videoImage, err = jr.repo.Media().Relate(*videoImage)
	if err != nil {
		return errs.BuildError(err, "could not create video image relation")
	}

	vidUpdate := dto.MediaOverviewDTO{
		Id:          video.Media.ID,
		ThumbnailId: image.MediaID,
	}
	jr.wsVideoUpdate(vidUpdate)

	return nil
}
