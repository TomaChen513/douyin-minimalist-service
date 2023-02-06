package lib

import (
	"bytes"
	"fmt"
	"io"
	"os"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func GetSnapshot(videoName, imageName string, frameNum int) (f io.Reader) {
	buf := bytes.NewBuffer(nil)
	ffmpeg.Input("https://douyinmiaomiao.oss-cn-hangzhou.aliyuncs.com/"+videoName).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()

	return buf

}
