package model

import (
	"log"
	"time"

	"github.com/RaymondCode/simple-demo/model/mysql"
)

type Message struct {
	Id      int64 `gorm:"primarykey" json:"id,omitempty"`
	UserId  int64 `json:"user_id,omitempty"`
	ToUserId int64 `json:"to_user_id,omitempty"`
	Content string
	CreateTime time.Time
}

// 表名映射
func (message Message) TableName() string{
	return "messages"
}

// 执行后向数据库插入一条信息数据
func InsertMessage(uId,tuId int64,content string,createTime time.Time)bool{
	message:=Message{UserId: uId,ToUserId: tuId,Content: content,CreateTime: createTime}
	if err:=mysql.DB.Create(&message).Error;err!=nil {
		log.Println(err.Error())
		return false
	}
	return true
}