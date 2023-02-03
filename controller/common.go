package controller

import "github.com/RaymondCode/simple-demo/service"

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type VideoListResponse struct {
	Response
	VideoList []service.Favor `json:"video_list"`
}

type VideoListResponse_publish struct {
	Response
	VideoList []service.Video `json:"video_list"`
}
