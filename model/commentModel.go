package model

import (
	"log"
	"time"

	"github.com/RaymondCode/simple-demo/model/mysql"
)

type TableComment struct {
	Id         int64
	UserId     int64
	VideoId    int64
	Content    string
	CreateDate time.Time
}

func (tableComment TableComment) TableName() string {
	return "comment"
}

func DeleteComment(commentId int64) bool {
	if err := mysql.DB.Delete(&TableComment{}, commentId).Error; err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

func InsertComment(tableComment TableComment) bool {
	if err := mysql.DB.Create(&tableComment).Error; err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

func GetCommentListDecByTime(videoId int64) ([]TableComment, error) {
	var tableComments []TableComment
	if err := mysql.DB.Where("video_id=?", videoId).
		Order("create_date desc").
		Find(&tableComments).
		Error; err != nil {
			return tableComments,err
	}
	return tableComments,nil
}
