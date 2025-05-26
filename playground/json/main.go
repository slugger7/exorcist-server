package main

import (
	"encoding/json"
	"fmt"

	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/ffmpeg"
)

// https://www.digitalocean.com/community/tutorials/how-to-use-json-in-go#parsing-json-using-a-map

func main() {
	jsonData := `{
    "streams": [
        {
            "index": 0,
            "codec_name": "h264",
            "codec_long_name": "H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10",
            "profile": "Main",
            "codec_type": "video",
            "codec_tag_string": "avc1",
            "codec_tag": "0x31637661",
            "width": 853,
            "height": 480,
            "coded_width": 853,
            "coded_height": 480,
            "closed_captions": 0,
            "film_grain": 0,
            "has_b_frames": 1,
            "pix_fmt": "yuv420p",
            "level": 30,
            "color_range": "tv",
            "color_space": "bt709",
            "color_transfer": "bt709",
            "color_primaries": "bt709",
            "chroma_location": "topleft",
            "field_order": "progressive",
            "refs": 1,
            "is_avc": "true",
            "nal_length_size": "4",
            "id": "0x1",
            "r_frame_rate": "24/1",
            "avg_frame_rate": "24/1",
            "time_base": "1/2400",
            "start_pts": 0,
            "start_time": "0.000000",
            "duration_ts": 1431500,
            "duration": "596.458333",
            "bit_rate": "2899884",
            "bits_per_raw_sample": "8",
            "nb_frames": "14315",
            "extradata_size": 36,
            "disposition": {
                "default": 1,
                "dub": 0,
                "original": 0,
                "comment": 0,
                "lyrics": 0,
                "karaoke": 0,
                "forced": 0,
                "hearing_impaired": 0,
                "visual_impaired": 0,
                "clean_effects": 0,
                "attached_pic": 0,
                "timed_thumbnails": 0,
                "non_diegetic": 0,
                "captions": 0,
                "descriptions": 0,
                "metadata": 0,
                "dependent": 0,
                "still_image": 0,
                "multilayer": 0
            },
            "tags": {
                "creation_time": "2008-05-27T18:32:32.000000Z",
                "language": "eng",
                "handler_name": "Apple Video Media Handler",
                "vendor_id": "appl",
                "encoder": "H.264"
            },
            "side_data_list": [
                {
                    "side_data_type": "Display Matrix",
                    "displaymatrix": "\n00000000:       116551           0           0\n00000001:            0      116599           0\n00000002:            0           0  1073741824\n",
                    "rotation": 0
                }
            ]
        },
        {
            "index": 1,
            "codec_type": "data",
            "codec_tag_string": "tmcd",
            "codec_tag": "0x64636d74",
            "id": "0x2",
            "r_frame_rate": "0/0",
            "avg_frame_rate": "2400/100",
            "time_base": "1/2400",
            "start_pts": 0,
            "start_time": "0.000000",
            "duration_ts": 1431500,
            "duration": "596.458333",
            "nb_frames": "1",
            "extradata_size": 22,
            "disposition": {
                "default": 1,
                "dub": 0,
                "original": 0,
                "comment": 0,
                "lyrics": 0,
                "karaoke": 0,
                "forced": 0,
                "hearing_impaired": 0,
                "visual_impaired": 0,
                "clean_effects": 0,
                "attached_pic": 0,
                "timed_thumbnails": 0,
                "non_diegetic": 0,
                "captions": 0,
                "descriptions": 0,
                "metadata": 0,
                "dependent": 0,
                "still_image": 0,
                "multilayer": 0
            },
            "tags": {
                "creation_time": "2008-05-27T18:32:32.000000Z",
                "language": "eng",
                "handler_name": "Time Code Media Handler",
                "timecode": "00:00:00:00"
            }
        },
        {
            "index": 2,
            "codec_name": "aac",
            "codec_long_name": "AAC (Advanced Audio Coding)",
            "profile": "LC",
            "codec_type": "audio",
            "codec_tag_string": "mp4a",
            "codec_tag": "0x6134706d",
            "sample_fmt": "fltp",
            "sample_rate": "48000",
            "channels": 6,
            "channel_layout": "5.1",
            "bits_per_sample": 0,
            "initial_padding": 0,
            "id": "0x4",
            "r_frame_rate": "0/0",
            "avg_frame_rate": "0/0",
            "time_base": "1/48000",
            "start_pts": 0,
            "start_time": "0.000000",
            "duration_ts": 28630160,
            "duration": "596.461667",
            "bit_rate": "437605",
            "nb_frames": "27960",
            "extradata_size": 2,
            "disposition": {
                "default": 1,
                "dub": 0,
                "original": 0,
                "comment": 0,
                "lyrics": 0,
                "karaoke": 0,
                "forced": 0,
                "hearing_impaired": 0,
                "visual_impaired": 0,
                "clean_effects": 0,
                "attached_pic": 0,
                "timed_thumbnails": 0,
                "non_diegetic": 0,
                "captions": 0,
                "descriptions": 0,
                "metadata": 0,
                "dependent": 0,
                "still_image": 0,
                "multilayer": 0
            },
            "tags": {
                "creation_time": "2008-05-27T18:32:32.000000Z",
                "language": "eng",
                "handler_name": "Apple Sound Media Handler",
                "vendor_id": "[0][0][0][0]"
            }
        }
    ],
    "format": {
        "filename": "./big_buck_bunny_480p_h264.mov",
        "nb_streams": 3,
        "nb_programs": 0,
        "nb_stream_groups": 0,
        "format_name": "mov,mp4,m4a,3gp,3g2,mj2",
        "format_long_name": "QuickTime / MOV",
        "start_time": "0.000000",
        "duration": "596.461667",
        "size": "249229883",
        "bit_rate": "3342778",
        "probe_score": 100,
        "tags": {
            "major_brand": "qt  ",
            "minor_version": "537199360",
            "compatible_brands": "qt  ",
            "creation_time": "2008-05-27T18:32:32.000000Z",
            "com.apple.quicktime.player.movie.audio.gain": "1.000000",
            "com.apple.quicktime.player.movie.audio.treble": "0.000000",
            "com.apple.quicktime.player.movie.audio.bass": "0.000000",
            "com.apple.quicktime.player.movie.audio.balance": "0.000000",
            "com.apple.quicktime.player.movie.audio.pitchshift": "0.000000",
            "com.apple.quicktime.player.movie.audio.mute": "",
            "com.apple.quicktime.player.movie.visual.brightness": "0.000000",
            "com.apple.quicktime.player.movie.visual.color": "1.000000",
            "com.apple.quicktime.player.movie.visual.tint": "0.000000",
            "com.apple.quicktime.player.movie.visual.contrast": "1.000000",
            "com.apple.quicktime.player.version": "7.4.1 (14)",
            "com.apple.quicktime.version": "7.4.1 (14) 0x7418000 (Mac OS X, 10.5.2, 9C31)",
            "timecode": "00:00:00:00"
        }
    }
}`

	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		fmt.Printf("could not unmarshal json: %s\n", err)

		return
	}

	var structdata *ffmpeg.Probe

	err = json.Unmarshal([]byte(jsonData), &structdata)
	errs.PanicError(err)

	width, height, err := ffmpeg.GetDimensions(structdata.Streams)
	errs.PanicError(err)

	fmt.Printf("Width & Height: %v %v\n", width, height)
}
