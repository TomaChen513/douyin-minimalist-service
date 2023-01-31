package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// Publish /publish/action/
func Publish(c *gin.Context) {
	// 鉴权token
	user_token := c.Query("token")
	_, flag := service.ParseToken(user_token)
	if !flag {
		c.JSON(http.StatusUnauthorized, Response{
			StatusCode: -1,
			StatusMsg:  "Token Error",
		})
		return
	}

	data, err := c.FormFile("data")
	if err != nil {
		log.Printf("获取视频流失败:%v", err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	log.Printf("获取到用户id:%v\n", userId)
	title := c.PostForm("title")
	log.Printf("获取到视频title:%v\n", title)

	file, err := data.Open()
	if err != nil {
		log.Printf("方法data.Open() 失败%v", err)
		return
	}
	log.Printf("方法data.Open() 成功")
	defer file.Close()

	videoName := uuid.NewV4().String()
	log.Printf("生成视频名称%v", videoName)

	err = model.VideoOss(file, videoName)
	if err != nil {
		log.Printf("方法videoService.Publish(data, userId) 失败：%v", err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	log.Printf("方法videoService.Publish(data, userId) 成功")

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  "uploaded successfully",
	})
}

// PublishList /publish/list/
func PublishList(c *gin.Context) {
	user_Id, _ := c.GetQuery("user_id")
	userId, _ := strconv.ParseInt(user_Id, 10, 64)
	log.Printf("获取到用户id:%v\n", userId)

	curId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	log.Printf("获取到当前用户id:%v\n", curId)

	// videoService := GetVideo()
	// list, err := videoService.List(userId, curId)
	// if err != nil {
	// 	log.Printf("调用videoService.List(%v)出现错误：%v\n", userId, err)
	// 	c.JSON(http.StatusOK, VideoListResponse{
	// 		Response: Response{StatusCode: 1, StatusMsg: "获取视频列表失败"},
	// 	})
	// 	return
	// }
	// log.Printf("调用videoService.List(%v)成功", userId)
	// c.JSON(http.StatusOK, VideoListResponse{
	// 	Response:  Response{StatusCode: 0},
	// 	VideoList: list,
	// })
}