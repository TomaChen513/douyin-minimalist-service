package controller

import (
	"net/http"
	"strconv"

	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
)

// mysql纯享版
// Post 发送消息操作 /douyin/message/action/
func MessageAction(c *gin.Context) {
	// 验证actionType
	actionType := c.Query("action_type")
	if actionType != "1" {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "actionType错误!"})
		return
	}
	// 根据token获取用户id
	user_id, _ := c.Get("userId")
	toUserId := c.Query("to_user_id")
	// 转换为int类型
	uId := user_id.(int64)
	tuId, _ := strconv.ParseInt(toUserId, 10, 64)
	// 获取信息内容
	content := c.Query("content")

	var msi service.MessageServiceImpl
	// 发送消息
	if ok := msi.MessageSend(uId, tuId, content); ok {
		c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "发送成功！"})
		return
	}
	c.JSON(http.StatusInternalServerError, Response{StatusCode: 1, StatusMsg: "发送失败！"})
}

// mysql纯享版
// 获得聊天记录 /douyin/message/chat/
func MessageChat(c *gin.Context) {
	// 根据token获取用户id
	user_id, _ := c.Get("userId")
	toUserId := c.Query("to_user_id")
	// 转换为int类型
	uId := user_id.(int64)
	tuId, _ := strconv.ParseInt(toUserId, 10, 64)

	print(uId, tuId)
	var msi service.MessageServiceImpl

	// 获取所有信息
	messages, err := msi.GetMessage(uId, tuId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{StatusCode: 1, StatusMsg: "获取消息失败！!"})
		return
	}

	c.JSON(http.StatusOK, MessageListResponse{
		Response:    Response{StatusCode: 0},
		MessageList: messages,
	})
}
