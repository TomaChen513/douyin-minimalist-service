package controller

import (
	"net/http"
	"strconv"

	"github.com/RaymondCode/simple-demo/lib"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
)

// var tempChat = map[string][]Message{
// 	"1_1": []Message{
// 		{Id: 1,
// 			Content:    "hello",
// 			CreateTime: "2023.1.16"},
// 	},
// }

// var messageIdSequence = int64(1)

type ChatResponse struct {
	Response
	MessageList []service.Message `json:"message_list"`
}

// 发送消息 POST /douyin/message/action/
func MessageAction(c *gin.Context) {
	token := c.Query("token")
	toUserId := c.Query("to_user_id")
	content := c.Query("content")
	actionType := c.Query("action_type")

	// 验证token
	tId, err := lib.GetKey(token)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "token信息有误！"})
		return
	}
	uId, _ := strconv.ParseInt(tId, 10, 64)
	tUId, _ := strconv.ParseInt(toUserId, 10, 64)

	if actionType != "1" {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "actionType有误！"})
		return
	}

	// 发送消息服务
	msi := service.MessageServiceImpl{}

	if !msi.SendMessage(uId,tUId,content) {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "发送消息有误！"})
		return
	}
	c.JSON(http.StatusOK, Response{StatusCode: 0})

	// if user, exist := usersLoginInfo[token]; exist {
	// 	userIdB, _ := strconv.Atoi(toUserId)
	// 	chatKey := genChatKey(user.Id, int64(userIdB))

	// 	atomic.AddInt64(&messageIdSequence, 1)
	// 	curMessage := Message{
	// 		Id:         messageIdSequence,
	// 		Content:    content,
	// 		CreateTime: time.Now().Format(time.Kitchen),
	// 	}

	// 	if messages, exist := tempChat[chatKey]; exist {
	// 		tempChat[chatKey] = append(messages, curMessage)
	// 	} else {
	// 		tempChat[chatKey] = []Message{curMessage}
	// 	}
	// 	c.JSON(http.StatusOK, Response{StatusCode: 0})
	// } else {
	// 	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	// }
}

// 聊天记录 GET /douyin/message/chat/
func MessageChat(c *gin.Context) {
	token := c.Query("token")
	toUserId := c.Query("to_user_id")

	// 验证token
	tId, err := lib.GetKey(token)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "token信息有误！"})
		return
	}
	uId, _ := strconv.ParseInt(tId, 10, 64)
	tUId, _ := strconv.ParseInt(toUserId, 10, 64)

	// 发送消息服务
	msi := service.MessageServiceImpl{}

	messageList,err:=msi.SelectMessageByUserId(uId,tUId)

	if err!=nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
	c.JSON(http.StatusOK, ChatResponse{Response: Response{StatusCode: 0}, MessageList: messageList})
	// if user, exist := usersLoginInfo[token]; exist {
	// 	userIdB, _ := strconv.Atoi(toUserId)
	// 	chatKey := genChatKey(user.Id, int64(userIdB))

	// 	c.JSON(http.StatusOK, ChatResponse{Response: Response{StatusCode: 0}, MessageList: tempChat[chatKey]})
	// } else {
	// 	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	// }
}

// func genChatKey(userIdA int64, userIdB int64) string {
// 	if userIdA > userIdB {
// 		return fmt.Sprintf("%d_%d", userIdB, userIdA)
// 	}
// 	return fmt.Sprintf("%d_%d", userIdA, userIdB)
// }
