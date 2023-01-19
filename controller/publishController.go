package controller

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/RaymondCode/simple-demo/lib"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
)

type VideoListResponse struct {
	Response
	VideoList []service.Video `json:"video_list"`
}

// 投稿接口  POST /douyin/publish/action/
func Publish(c *gin.Context) {
	// 验证接口，上传文件，上传数据库记录
	token := c.PostForm("token")


	vsi := service.VideoServiceImpl{}

	// 验证token
	tId, err := lib.GetKey(token)
	if err!=nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "token信息有误！"},
		})
		return
	}
	tokenId, _ := strconv.ParseInt(tId, 10, 64)

	// 上传文件
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	filename := filepath.Base(data.Filename)
	user,_:=vsi.GetUserById(tokenId)
	finalName := fmt.Sprintf("%d_%s", user.Id, filename)
	saveFile := filepath.Join("./public/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	// 添加数据进入数据库
	video:=service.Video{
		Author: user,
		PlayUrl: "./public/"+finalName,
		CoverUrl: "./public/"+finalName,
		FavoriteCount: 0,
		CommentCount: 0,
		IsFavorite: false,
		Title: c.PostForm("title"),
	}
	if ok:=vsi.InsertVideo(video);!ok{
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  finalName + " 数据库插入失败",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// 获得用户发布信息 GET /douyin/publish/list/
func PublishList(c *gin.Context) {
	token := c.Query("token")
	userId := c.Query("user_id")
	id, _ := strconv.ParseInt(userId, 10, 64)

	vsi := service.VideoServiceImpl{}

	// 验证token
	tId, _ := lib.GetKey(token)
	tokenId, _ := strconv.ParseInt(tId, 10, 64)
	if id != tokenId {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户id与token信息不一致！"},
		})
		return
	}

	if _, err := model.GetUserById(id); err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户不存在"},
		})
	} else {
		videoList, _ := vsi.GetVideosByUser(id)
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 0,
			},
			VideoList: videoList,
		})
	}

}
