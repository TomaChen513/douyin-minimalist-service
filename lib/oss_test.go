package lib

import (
	"fmt"
	"testing"
)

func TestUploadOss(t *testing.T) {
	UploadOss("bear1.mp4","filehash")
	fmt.Println("success")
}

func TestDownLoadOss(t *testing.T){
	data:=DownloadOss("filehash",".mp4")
	fmt.Println(data)
}

func TestDeleteOss(t *testing.T){
	DeleteOss("filehash",".mp4")
}