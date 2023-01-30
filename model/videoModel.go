package model

import (
	"log"
	"time"

	"github.com/RaymondCode/simple-demo/model/mysql"
)

type Video struct {
	Id          int64  `gorm:"primarykey" json:"id,omitempty"`
	AuthorId    int64  `gorm:"type:int" json:"author_id"`
	PlayUrl     string `gorm:"type:varchar(255)" json:"play_url,omitempty"`
	CoverUrl    string `gorm:"type:varchar(255)" json:"cover_url,omitempty"`
	PublishTime time.Time
}

// TableName 修改表名映射
func (Video) TableName() string {
	return "videos"
}

func GetVideosById(videoId []int64) ([]Video, error) {
	var videos []Video
	if err := mysql.DB.Find(&videos, videoId).Error; err != nil {
		log.Println(err.Error())
		return []Video{}, err
	}
	return videos, nil
}
