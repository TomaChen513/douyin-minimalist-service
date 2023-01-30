package controller

import (
	"strconv"

	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
)

// 关注操作
func RelationAction(c *gin.Context) {
	//解析token
	token, ok := service.ParseToken(c.Query("token"))

	if !ok {
		c.JSON(200, gin.H{
			"StatusCode": -1,
			"StatusMsg":  "鉴权失败",
		})
	}

	followerId, err1 := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	cancel, err2 := strconv.ParseInt(c.Query("action_type"), 10, 64)

	//判断参数格式
	if err1 != nil || err2 != nil {
		if !ok {
			c.JSON(200, gin.H{
				"StatusCode": -1,
				"StatusMsg":  "用户id格式或行为格式错误",
			})
		}
	}

	//更新数据
	if ok = service.FollowAction(token.UserId, followerId, int8(cancel)); !ok {
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

func FollowList(c *gin.Context) {
	//解析token
	_, ok := service.ParseToken(c.Query("token"))

	if !ok {
		c.JSON(200, gin.H{
			"StatusCode": -1,
			"StatusMsg":  "鉴权失败",
		})
	}

	userId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)

	//判断参数格式
	if err != nil {
		c.JSON(200, gin.H{
			"StatusCode": -1,
			"StatusMsg":  "用户id格式错误",
			"user_list":  nil,
		})
	}

	users, ok := service.GetFollowList(userId)

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

func FollowerList(c *gin.Context) {
	//解析token
	_, ok := service.ParseToken(c.Query("token"))

	if !ok {
		c.JSON(200, gin.H{
			"StatusCode": -1,
			"StatusMsg":  "鉴权失败",
		})
	}

	userId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)

	//判断参数格式
	if err != nil {
		c.JSON(200, gin.H{
			"StatusCode": -1,
			"StatusMsg":  "用户id格式错误",
			"user_list":  nil,
		})
	}

	users, ok := service.GetFollowerList(userId)

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
