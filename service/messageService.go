package service

import (
	"log"
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
	// 根据双方id获得信息
	GetMessage(uId,tuId int64)
}

type MessageServiceImpl struct {
	
}

func (msi *MessageServiceImpl) MessageSend(uId,tuId int64,content string) bool{
	// 调用model函数
	createTime:=time.Now()
	return model.InsertMessage(uId,tuId,content,createTime)
}

func(msi *MessageServiceImpl)GetMessage(uId,tuId int64) ([]Message,error){
	messages,err:=model.SelectMessageByUserId(uId,tuId)
	serviceMessages:=make([]Message,0)
	if err!=nil {
		log.Println(err.Error())
		return serviceMessages,err
	}
	for _,m:=range messages{
		serviceMessages = append(serviceMessages, Message{
			Id: m.Id,
			Content: m.Content,
			CreateTime: m.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}
	return serviceMessages,nil
}