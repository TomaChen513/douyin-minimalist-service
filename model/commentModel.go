package model

import (
	"log"
	"time"

	"github.com/RaymondCode/simple-demo/model/mysql"
)

type Comment struct {
	Id         int64
	UserId     int64
	VideoId    int64
	CommentText    string
	CreateDate time.Time
	Cancel int64
}

func (comment Comment) TableName() string {
	return "comments"
}

func DeleteComment(commentId int64) bool {
	if err := mysql.DB.Delete(&Comment{}, commentId).Error; err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

func InsertComment(comment Comment) bool {
	if err := mysql.DB.Create(&comment).Error; err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

func GetCommentListDecByTime(videoId int64) ([]Comment, error) {
	var comments []Comment
	if err := mysql.DB.Where("video_id=?", videoId).
		Order("create_date desc").
		Find(&comments).
		Error; err != nil {
			return comments,err
	}
	return comments,nil
}