package model

import (
	"log"

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