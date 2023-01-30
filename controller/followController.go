package controller

import (
	"strconv"

	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
)

// 关注操作 userId表示当前用户Id, followerId, 表示关注对象
func RelationAction(c *gin.Context) {
	userId, err1 := strconv.ParseInt(c.GetString("userId"), 10, 64)
	followerId, err2 := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	cancel, err3 := strconv.ParseInt(c.Query("action_type"), 10, 64)

	//判断参数格式
	if err1 != nil || err2 != nil || err3 != nil {
		c.JSON(200, gin.H{
			"StatusCode": -1,
			"StatusMsg":  "参数格式错误",
		})
	}

	//更新数据
	if ok := service.FollowAction(userId, followerId, int8(cancel)); !ok {
		c.JSON(200, gin.H{
			"StatusCode": -1,
			"StatusMsg":  "更新数据失败",
		})
	}

	c.JSON(200, gin.H{
		"StatusCode": 0,
		"StatusMsg":  "成功",
	})
}

// 关注列表， curId当前登录用户id， userId查询对象
func FollowList(c *gin.Context) {
	curId, err1 := strconv.ParseInt(c.GetString("userId"), 10, 64)
	userId, err2 := strconv.ParseInt(c.Query("to_user_id"), 10, 64)

	//判断参数格式
	if err1 != nil || err2 != nil {
		c.JSON(200, gin.H{
			"StatusCode": -1,
			"StatusMsg":  "参数格式错误",
			"user_list":  nil,
		})
	}

	users, ok := service.GetFollowList(userId, curId)

	if !ok {
		c.JSON(200, gin.H{
			"StatusCode": -1,
			"StatusMsg":  "获取关注列表失败",
			"user_list":  nil,
		})
	}

	c.JSON(200, gin.H{
		"StatusCode": 0,
		"StatusMsg":  "获取关注列表成功",
		"user_list":  users,
	})
}

// 粉丝列表， curId当前登录用户id， userId查询对象
func FollowerList(c *gin.Context) {
	curId, err1 := strconv.ParseInt(c.GetString("userId"), 10, 64)
	userId, err2 := strconv.ParseInt(c.Query("to_user_id"), 10, 64)

	//判断参数格式
	if err1 != nil || err2 != nil {
		c.JSON(200, gin.H{
			"StatusCode": -1,
			"StatusMsg":  "参数格式错误",
			"user_list":  nil,
		})
	}

	users, ok := service.GetFollowerList(userId, curId)

	if !ok {
		c.JSON(200, gin.H{
			"StatusCode": -1,
			"StatusMsg":  "获取粉丝列表失败",
			"user_list":  nil,
		})
	}

	c.JSON(200, gin.H{
		"StatusCode": 0,
		"StatusMsg":  "获取粉丝列表成功",
		"user_list":  users,
	})
}
