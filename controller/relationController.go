package controller

import (
	"net/http"
	"strconv"

	"github.com/RaymondCode/simple-demo/lib"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
)

type UserListResponse struct {
	Response
	UserList []service.User `json:"user_list"`
}

// 关注操作  POST /douyin/relation/action/
func RelationAction(c *gin.Context) {
	token := c.Query("token")
	toUserId := c.Query("to_user_id")
	actionType := c.Query("action_type")
	tUId,_:=strconv.ParseInt(toUserId,10,64)

	// 验证token
	tId, err := lib.GetKey(token)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "token信息有误！"})
		return
	}
	uId, _ := strconv.ParseInt(tId, 10, 64)

	ssi := service.SocialServiceImpl{}

	// 执行关注或者取消关注操作
	if !ssi.AttentionAction(uId, tUId, actionType) {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "关注操作数据库操作失败！"})
		return
	}
	c.JSON(http.StatusOK, Response{StatusCode: 0})
}

// 关注列表 GET /douyin/relation/follow/list/
func FollowList(c *gin.Context) {
	token := c.Query("token")
	userId := c.Query("user_id")
	uId, _ := strconv.ParseInt(userId, 10, 64)

	ssi := service.SocialServiceImpl{}

	// 验证token
	_, err := lib.GetKey(token)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "token信息有误！"},
		})
		return
	}

	followList, err := ssi.GetFollowList(uId)

	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "数据库处理错误！"},
		})
		return
	} else {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 0,
			},
			UserList: followList,
		})
	}
}

// 粉丝列表  GET /douyin/relation/follower/list/
func FollowerList(c *gin.Context) {
	token := c.Query("token")
	userId := c.Query("user_id")
	uId, _ := strconv.ParseInt(userId, 10, 64)

	ssi := service.SocialServiceImpl{}

	// 验证token
	_, err := lib.GetKey(token)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "token信息有误！"},
		})
		return
	}

	followerList, err := ssi.GetFollowerList(uId)

	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "数据库处理错误！"},
		})
		return
	} else {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 0,
			},
			UserList: followerList,
		})
	}
}

// 好友列表 GET /douyin/relation/friend/list/
func FriendList(c *gin.Context) {
	token := c.Query("token")
	userId := c.Query("user_id")
	uId, _ := strconv.ParseInt(userId, 10, 64)

	ssi := service.SocialServiceImpl{}

	// 验证token
	_, err := lib.GetKey(token)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "token信息有误！"},
		})
		return
	}

	friendList, err := ssi.GetFollowList(uId)

	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "数据库处理错误！"},
		})
		return
	} else {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 0,
			},
			UserList: friendList,
		})
	}
}
