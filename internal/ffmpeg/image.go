package ffmpeg

import (
	"fmt"

	errs "github.com/slugger7/exorcist/internal/errors"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

const (
	ErrExtractingImage string = "error extracting image (%v) from video (%v) at (%v) with resolution %vx%v"
	ErrNegativeWidth   string = "width cannot be negative or zero: %v"
	ErrNegativeHeight  string = "height cannot be negative or zero: %v"
)

func ScaleWidthByHeight(currentHeight, currentWidth, wantedHeight int) int {
	return int(float32(currentWidth) / float32(currentHeight) * float32(wantedHeight))
}

func ScaleHeightByWidth(currentHeight, currentWidth, wantedWidth int) int {
	return int(float32(currentHeight) / float32(currentWidth) * float32(wantedWidth))
}

func ImageAt(vid string, time int, img string, width, height int) error {
	if width <= 0 {
		return fmt.Errorf(ErrNegativeWidth, width)
	}
	if height <= 0 {
		return fmt.Errorf(ErrNegativeHeight, height)
	}

	err := ffmpeg_go.Input(vid, ffmpeg_go.KwArgs{"ss": time}).
		Output(img, ffmpeg_go.KwArgs{"vframes": 1, "s": fmt.Sprintf("%vx%v", width, height)}).
		Run()

	if err != nil {
		return errs.BuildError(err, ErrExtractingImage, img, vid, time, width, height)
	}

	return nil
}
