package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	Response
	VideoList []service.Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

// GET 视频流接口  /douyin/feed/
func Feed(c *gin.Context) {
	// token := c.Query("token")
	latest_time:=c.Query("latest_time")
	latestTime,_:=strconv.ParseInt(latest_time,10,64)

	if latest_time=="0" {
		latestTime=time.Now().Unix()
	}

	vsi := service.VideoServiceImpl{}

	videoList,nextTime,err:=vsi.Feed(latestTime)

	if err!=nil {
		// 失败传回状态码1
		c.JSON(http.StatusOK, FeedResponse{
			Response:  Response{StatusCode: 1},
			VideoList: videoList,
			NextTime:  time.Now().Unix(),
		})
		return
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videoList,
		// NextTime:  time.Now().Unix(),
		NextTime: nextTime,
	})
}
