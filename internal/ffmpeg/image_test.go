package ffmpeg

import (
	"fmt"
	"os"
	"testing"

	"github.com/slugger7/exorcist/internal/assert"
)

func Test_ImageAt_NegativeWidth(t *testing.T) {
	width := -1

	err := ImageAt("", 0, "", width, 1)
	assert.ErrorNotNil(t, err)
	assert.Error(t, fmt.Errorf(ErrNegativeWidth, width), err)
}

func Test_ImageAt_NegativeHeight(t *testing.T) {
	height := -1

	err := ImageAt("", 0, "", 1, height)
	assert.ErrorNotNil(t, err)
	assert.Error(t, fmt.Errorf(ErrNegativeHeight, height), err)
}

func Test_ImageAt_Success(t *testing.T) {
	width, height, time := 20, 60, 3

	err := ImageAt(testVideoPath, time, testImagePath, width, height)
	assert.ErrorNil(t, err)
	assert.FileExists(t, testImagePath)

	os.Remove(testImagePath)
}
