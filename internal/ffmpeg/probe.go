package ffmpeg

import (
	"encoding/json"
	"errors"

	ffmpegGo "github.com/u2takey/ffmpeg-go"
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

func UnmarshalProbeData(probeData string) (*Probe, error) {
	var data *Probe
	err := json.Unmarshal([]byte(probeData), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func UnmarshalledProbe(path string) (*Probe, error) {
	probeData, err := ffmpegGo.Probe(path)
	if err != nil {
		return nil, err
	}

	data, err := UnmarshalProbeData(probeData)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetDimensions(streams []Stream) (width, height int, err error) {
	for _, v := range streams {
		if v.CodecType == "video" {
			return *v.Width, *v.Height, nil
		}
	}

	return 0, 0, errors.New("could not extract the height and with from the probe data streams")
}
