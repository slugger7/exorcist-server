package ffmpeg

import (
	"reflect"
	"strings"
	"testing"
)

func Test_GetDImensions_NoVideoCodecInStreams_shouldCreateError(t *testing.T) {
	sterams := []Stream{
		{
			CodecType: "not_video",
		},
	}

	expectedError := "could not extract the height and with from the probe data streams"

	_, _, err := GetDimensions(sterams)

	if err != nil {
		if err.Error() != expectedError {
			t.Errorf("Expected error to be '%v' but was '%v'", expectedError, err.Error())
		}
	} else {
		t.Error("Expected an error but none was thrown")
	}
}

func Test_GetDImensions_WithVideoCodec_shouldReturnHeightAndWidth_withNilError(t *testing.T) {
	width, height := 69, 420
	streams := []Stream{
		{
			CodecType: "video",
			Width:     &width,
			Height:    &height,
		},
	}

	actualWidth, actualHeight, err := GetDimensions(streams)
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

func Test_UnmarshalledProbe_WithBrokenTestVideoFile_ShouldThrowError(t *testing.T) {
	_, err := UnmarshalledProbe(brokenVideoPath)
	if err != nil {
		if !strings.Contains(err.Error(), "Invalid data found when processing input") {
			t.Errorf("Incorrect error was thrown: %v", err)
			return
		} else {
			return
		}
	}

	t.Errorf("No error was thrown")
}

func Test_UnmarshalProbeData_WithAWorkingFile_ShouldCreateCorrectStruct(t *testing.T) {
	actual, err := UnmarshalledProbe(testVideoPath)
	if err != nil {
		t.Errorf("Error was thrown %v", err)
		return
	}

	height, width := 270, 480
	expectedFormat := Format{
		Duration: "33.023333",
		Size:     "3889885",
	}
	expectedStream := Stream{
		Height:    &height,
		Width:     &width,
		CodecType: "video",
	}

	var actualVideoStream *Stream
	for _, v := range actual.Streams {
		if v.CodecType == expectedStream.CodecType {
			actualVideoStream = &v
		}
	}
	if actualVideoStream == nil {
		t.Error("Could not find a video stream")
		return
	}

	if !reflect.DeepEqual(*actual.Format, expectedFormat) {
		t.Error("Actual format does not match expected format")
	}
	if !reflect.DeepEqual(*actualVideoStream, expectedStream) {
		t.Error("Actual video stream does not match expected video stream")
	}
}

func Test_UnmarshalProbeData_WithInvalidJson_ShouldThrowError(t *testing.T) {
	_, err := UnmarshalProbeData("this is not json")
	if err != nil {
		if !strings.Contains(err.Error(), "invalid character 'h' in literal true (expecting 'r')") {
			t.Errorf("Incorrect error was thrown: %v", err)
			return
		} else {
			return
		}
	}
	t.Error("No error was thrown")
}

func Test_UnmarshallProbeData_WithValidJson_ShouldParseJson(t *testing.T) {
	jsonData := `{
		"format": {
			"duration": "66.6",
			"size": "666"
		},
		"streams": [
			{
				"codec_type": "video",
				"height": 69,
				"width": 420
			}
		]
	}`
	expectedHeight, expectedWidth := 69, 420
	expectedProbeData := Probe{
		Format: &Format{
			Duration: "66.6",
			Size:     "666",
		},
		Streams: []Stream{
			{
				CodecType: "video",
				Height:    &expectedHeight,
				Width:     &expectedWidth,
			},
		},
	}
	data, err := UnmarshalProbeData(jsonData)
	if err != nil {
		t.Errorf("Could not unmarshal json data: %v", err)
	}

	if !reflect.DeepEqual(*data, expectedProbeData) {
		t.Errorf("Expected data differed from actual data")
	}
}
