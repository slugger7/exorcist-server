package main

import (
	"fmt"

	. "github.com/slugger7/exorcist/internal/errors"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func main() {
	a, err := ffmpeg.Probe("<insert video path here>")
	CheckError(err)

	fmt.Println(a)
}
