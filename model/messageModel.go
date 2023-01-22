package model

import (
	"log"
	"time"

	"github.com/RaymondCode/simple-demo/model/mysql"
)

type TableMessage struct {
	Id         int64  `json:"id,omitempty"`
	Content    string `json:"content,omitempty"`
	CreateTime time.Time `json:"create_time,omitempty"`
	UserId int64
	ToUserId int64
}


func (tableMessage TableMessage) TableName()string{
	return "message"
}

func InsertMessage(tableMessage TableMessage) bool{
	tableMessage.CreateTime=time.Now()
	if err := mysql.DB.Create(&tableMessage).Error; err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

func SelectMessagesByUserId(userId,toUserId int64) ([]TableMessage,error){
	var tableMessages []TableMessage
	if err:=mysql.DB.Where("user_id=?",userId).
	Where("to_user_id=?",toUserId).
	Find(&tableMessages).Error;err!=nil {
		log.Println(err.Error())
		return tableMessages,err
	}
	return tableMessages,nil
}