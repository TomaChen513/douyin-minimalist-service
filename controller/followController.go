package controller

import (
	"fmt"
	"strconv"

	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
)

// 关注操作 userId表示当前用户Id, followerId, 表示关注对象
// 1. 没有考虑to_user_id不存在的情况(已修改)  2. 可以直接取消关注，即用户没有关注的时候，都可以取关（已修改）
func RelationAction(c *gin.Context) {
	user_id, _ := c.Get("userId")
	userId := user_id.(int64)
	followerId, err1 := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	cancel, err2 := strconv.ParseInt(c.Query("action_type"), 10, 64)

	//判断参数格式
	if err1 != nil || err2 != nil {
		c.JSON(200, gin.H{
			"StatusCode": -1,
			"StatusMsg":  "参数格式错误",
		})
		return
	}

	usi := service.UserServiceImpl{}
	_, err3 := usi.GetUserById(followerId)

	if err3 != nil {
		c.JSON(200, gin.H{
			"StatusCode": -1,
			"StatusMsg":  "to_user_id不存在",
		})
		return
	}

	//更新数据
	if ok := service.FollowAction(userId, followerId, int8(cancel)); !ok {
		c.JSON(200, gin.H{
			"StatusCode": -1,
			"StatusMsg":  "更新数据失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"StatusCode": 0,
		"StatusMsg":  "成功",
	})
}

// 关注列表， curId当前登录用户id， userId查询对象
func FollowList(c *gin.Context) {
	cur_id, ok := c.Get("userId")
	curId := cur_id.(int64)

	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)

	println(curId, userId)

	//判断参数格式
	if !ok || err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"StatusCode": -1,
			"StatusMsg":  "参数格式错误",
			"user_list":  nil,
		})
		return
	}

	users, ok := service.GetFollowList(userId, curId)

	if !ok {
		c.JSON(200, gin.H{
			"StatusCode": -1,
			"StatusMsg":  "获取关注列表失败",
			"user_list":  nil,
		})
		return
	}

	c.JSON(200, gin.H{
		"StatusCode": 0,
		"StatusMsg":  "获取关注列表成功",
		"user_list":  users,
	})
	// c.JSON(200, gin.H{
	// 	"test": "11",
	// })
}

// 粉丝列表， curId当前登录用户id， userId查询对象
func FollowerList(c *gin.Context) {
	cur_id, ok := c.Get("userId")
	curId := cur_id.(int64)
	userId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)

	//判断参数格式
	if !ok || err != nil {
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
		return
	}

	c.JSON(200, gin.H{
		"StatusCode": 0,
		"StatusMsg":  "获取粉丝列表成功",
		"user_list":  users,
	})
}
