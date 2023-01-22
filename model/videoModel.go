package model

import (
	"log"
	"time"

	"github.com/RaymondCode/simple-demo/model/mysql"
)


type TableVideo struct {
	Id            int64  `gorm:"primarykey" json:"id,omitempty"`
	Author        int64   `gorm:"type:int" json:"author"`
	PlayUrl       string `gorm:"type:varchar(255)" json:"play_url,omitempty"`
	CoverUrl      string `gorm:"type:varchar(255)" json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
	Title string `gorm:"type:varchar(55)" json:"title,omitempty"`
	PublishTime time.Time
}

// TableName 修改表名映射
func (tableVideo TableVideo) TableName() string {
	return "video"
}


// 根据用户id获取所有投稿信息
func GetVideosByUserId(userId int64) ([]TableVideo,error){
	videos := []TableVideo{}
	if err:=mysql.DB.Where("author = ? ", userId).Find(&videos).Error;err!=nil{
		log.Println(err.Error())
		return videos,err
	}
	return videos, nil
}

// 插入视频信息
func InsertVideo(tableVideo TableVideo)bool{
	if err := mysql.DB.Create(&tableVideo).Error; err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

// 根据主键数组查找video
func SelectVideosByPriArr(ids []int64) []TableVideo{
	var tableVideo []TableVideo
	mysql.DB.Find(&tableVideo,ids)
	return tableVideo
}

func GetVideosByLatestTime(lastTime time.Time) ([]TableVideo,error){
	var tableVideo []TableVideo
	result := mysql.DB.Where("publish_time<?", lastTime).Order("publish_time desc").Limit(30).Find(&tableVideo)
	if result.Error != nil {
		return tableVideo, result.Error
	}
	return tableVideo, nil
}