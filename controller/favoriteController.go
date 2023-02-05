package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
)

// 有重复点赞的bug， 在点赞前先查询数据库中是否已经存在记录
// 点赞操作  POST /douyin/favorite/action/
func FavoriteAction(c *gin.Context) {
	userId := c.GetString("userId")
	videoId := c.Query("video_id")
	actionType := c.Query("action_type")
	vId, _ := strconv.ParseInt(videoId, 10, 64)
	uId, _ := strconv.ParseInt(userId, 10, 64)

	fvsi := service.FavorServiceImpl{}

	// 执行点赞或者取消点赞操作
	if !fvsi.FavoriteAction(uId, vId, actionType) {
		log.Printf("service.FavoriteAction(uId, vId, actionType) 失败")
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "点赞操作失败！"})
		return
	}
	log.Printf("service.FavoriteAction(uId, vId, actionType) 成功")
	c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "点赞操作成功！"})
}

// 不知道为啥一直返回空数组， 可能我测试有问题
// 喜欢列表 GET /douyin/favorite/list/
func FavoriteList(c *gin.Context) {
	// user_id：待查询的用户id，token只是为了验证是否有权限
	userId := c.Query("user_id")
	uId, _ := strconv.ParseInt(userId, 10, 64)

	fvsi := service.FavorServiceImpl{}

	// 获得喜欢列表
	favorVideoList, err := fvsi.GetFavouriteList(uId)

	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1,
				StatusMsg: "GetFavouriteList(userId int64) ([]Favor, error)调用失败"},
		})
		return
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "GetFavouriteList(userId int64) ([]Favor, error)调用成功"},
		VideoList: favorVideoList,
	})
}
