package ffmpeg

import (
	"errors"
)

type Stream struct {
	Height    *int   `json:"height"`
	Width     *int   `json:"width"`
	CodecType string `json:"codec_type"`
}

type Format struct {
	Duration string `json:"duration"`
	Size     string `json:"size"`
}

type Probe struct {
	Format  *Format  `json:"format"`
	Streams []Stream `json:"streams"`
}

func GetDimensions(streams []Stream) (width, height int, err error) {
	for _, v := range streams {
		if v.CodecType == "video" {
			return *v.Width, *v.Height, nil
		}
	}

	return 0, 0, errors.New("could not extract the height and with from the probe data streams")
}
