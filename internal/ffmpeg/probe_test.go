package ffmpeg_test

import (
	"testing"

	"github.com/slugger7/exorcist/internal/ffmpeg"
)

func TestNoVideoCodecInStreams_shouldCreateError(t *testing.T) {
	sterams := []ffmpeg.Stream{
		{
			CodecType: "not_video",
		},
	}

	expectedError := "could not extract the height and with from the probe data streams"

	_, _, err := ffmpeg.GetDimensions(sterams)

	if err != nil {
		if err.Error() != expectedError {
			t.Errorf("Expected error to be '%v' but was '%v'", expectedError, err.Error())
		}
	} else {
		t.Error("Expected an error but none was thrown")
	}
}

func TestWithVideoCodec_shouldReturnHeightAndWidth_withNilError(t *testing.T) {
	width, height := 69, 420
	streams := []ffmpeg.Stream{
		{
			CodecType: "video",
			Width:     &width,
			Height:    &height,
		},
	}

	actualWidth, actualHeight, err := ffmpeg.GetDimensions(streams)
	if err != nil {
		t.Errorf("Could not extract height and width from streams with error %v", err)
	}
	if width != actualWidth {
		t.Errorf("Actual width (%v) does not match expected width (%v)", actualWidth, width)
	}
	if height != actualHeight {
		t.Errorf("Actual height (%v) does not match expected height (%v)", actualHeight, height)
	}
}
