package service

import (
	"time"

	"github.com/RaymondCode/simple-demo/model"
)

type Message struct {
	Id         int64  `json:"id,omitempty"`
	Content    string `json:"content,omitempty"`
	CreateTime string `json:"create_time"`
}

type MessageService interface {
	// 发送消息
	MessageSend(uId,tuId int64,content string) bool
}

type MessageServiceImpl struct {
	
}

func (msi *MessageServiceImpl) MessageSend(uId,tuId int64,content string) bool{
	// 调用model函数
	createTime:=time.Now()
	return model.InsertMessage(uId,tuId,content,createTime)
}