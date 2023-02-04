package controller

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

type FeedResponse struct {
	Response
	VideoList []service.Video `json:"video_list"`
	NextTime  int64           `json:"next_time"`
}

// Publish /publish/action/
func Publish(c *gin.Context) {

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

	//封面图片名称（暂时）
	imageName := videoName

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

	//组装并持久化
	_ = model.Save(videoName, imageName, userId, title)

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  "uploaded successfully",
	})
}

// Feed /feed/
func Feed(c *gin.Context) {
	inputTime := c.Query("latest_time")
	log.Printf("传入的时间" + inputTime)
	var lastTime time.Time
	if inputTime != "0" {
		me, _ := strconv.ParseInt(inputTime, 10, 64)
		lastTime = time.Unix(me, 0)
	} else {
		lastTime = time.Now()
	}
	log.Printf("获取到时间戳%v", lastTime)
	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	log.Printf("获取到用户id:%v\n", userId)
	videoService := GetVideo()
	feed, nextTime, err := videoService.Feed(lastTime, userId)
	if err != nil {
		log.Printf("方法videoService.Feed(lastTime, userId) 失败：%v", err)
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: 1, StatusMsg: "获取视频流失败"},
		})
		return
	}
	log.Printf("方法videoService.Feed(lastTime, userId) 成功")
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: feed,
		NextTime:  nextTime.Unix(),
	})
}

// PublishList /publish/list/
func PublishList(c *gin.Context) {
	user_Id, _ := c.GetQuery("user_id")
	userId, _ := strconv.ParseInt(user_Id, 10, 64) //被查询目标的id
	log.Printf("获取到待查询目标用户id:%v\n", userId)

	curId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64) //现在登录账号的人的id
	log.Printf("获取到当前用户id:%v\n", curId)

	videoService := GetVideo()
	list, err := videoService.List(userId, curId)
	if err != nil {
		log.Printf("调用videoService.List(%v)出现错误：%v\n", userId, err)
		c.JSON(http.StatusOK, VideoListResponse_publish{
			Response: Response{StatusCode: 1, StatusMsg: "获取视频列表失败"},
		})
		return
	}
	log.Printf("调用videoService.List(%v)成功", userId)
	c.JSON(http.StatusOK, VideoListResponse_publish{
		Response:  Response{StatusCode: 0},
		VideoList: list,
	})
}

// GetVideo 拼装videoService
func GetVideo() service.VideoServiceImpl {
	var userService service.UserServiceImpl
	var followService service.FollowServiceImp
	var videoService service.VideoServiceImpl
	var likeService service.FavorServiceImpl
	var commentService service.CommentServiceImpl

	userService.FollowService = &followService
	userService.FavorService = &likeService
	followService.UserService = &userService
	likeService.VideoService = &videoService
	commentService.UserService = &userService
	videoService.CommentService = &commentService
	videoService.FavorService = &likeService
	videoService.UserService = &userService

	return videoService
}
