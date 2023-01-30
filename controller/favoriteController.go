package controller

import (
	"net/http"
	"strconv"

	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
)

// 点赞操作  POST /douyin/favorite/action/
func FavoriteAction(c *gin.Context) {
	// 使用mysql进行操作，性能极差，建议优化
	userId := c.GetString("userId")
	videoId := c.Query("video_id")
	actionType := c.Query("action_type")
	vId, _ := strconv.ParseInt(videoId, 10, 64)
	uId,_:=strconv.ParseInt(userId,10,64)

	fvsi := service.FavorServiceImpl{}

	// 执行点赞或者取消点赞操作
	if !fvsi.FavoriteAction(uId, vId, actionType) {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "点赞操作数据库操作失败！"})
		return
	}
	c.JSON(http.StatusOK, Response{StatusCode: 0})
}

// 喜欢列表 GET /douyin/favorite/list/
func FavoriteList(c *gin.Context) {
	userId := c.Query("user_id")
	uId, _ := strconv.ParseInt(userId, 10, 64)

	fvsi := service.FavorServiceImpl{}

	favorVideoList, err := fvsi.GetFavouriteList(uId)

	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "数据库处理错误！"},
		})
		return
	} else {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 0,
			},
			VideoList: favorVideoList,
		})
	}
}