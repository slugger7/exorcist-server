package main

import "github.com/slugger7/exorcist/internal/ffmpeg"

// ffmpeg -ss 00:00:04 -i $PWD/internal/ffmpeg/test_data/working_video.mp4 -frames:v 1 $PWD/.temp/screenshot.png
// https://www.bannerbear.com/blog/how-to-extract-images-from-a-video-using-ffmpeg/
func main() {
	vid := "./internal/ffmpeg/test_data/working_video.mp4"
	img := "./.temp/img.png"

	err := ffmpeg.ImageAt(vid, 32, img, 30, 20)
	if err != nil {
		panic(err)
	}
}
