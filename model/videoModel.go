package model

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/RaymondCode/simple-demo/lib"
	"github.com/RaymondCode/simple-demo/model/mysql"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type Video struct {
	Id          int64  `gorm:"primarykey" json:"id,omitempty"`
	AuthorId    int64  `gorm:"type:int" json:"author_id"`
	PlayUrl     string `gorm:"type:varchar(255)" json:"play_url,omitempty"`
	CoverUrl    string `gorm:"type:varchar(255)" json:"cover_url,omitempty"`
	Title       string `json:"title"` //视频名
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

// 把视频上传到oss
func VideoOss(file io.Reader, videoName string) error {
	conf := lib.LoadServerConfig()
	// 创建OSSClient实例。
	client, err := oss.New(conf.Endpoint, conf.AccessKeyID, conf.AccessKeySecret)
	if err != nil {
		fmt.Println("创建实例Error:", err)
		return err
	}
	// 获取存储空间。
	bucket, err := client.Bucket(conf.BucketName)
	if err != nil {
		fmt.Println("获取存储空间Error:", err)
		return err
	}

	// 上传文件。
	err = bucket.PutObject(videoName, file)
	if err != nil {
		fmt.Println("文件上传Error:", err)
		return err
	}

	fmt.Println("上传文件成功！")
	return nil
}

// GetVideosByAuthorId
// 根据作者的id来查询对应数据库数据，并TableVideo返回切片
func GetVideosByAuthorId(authorId int64) ([]Video, error) {
	//建立结果集接收
	var data []Video
	//初始化db
	//Init()
	result := mysql.DB.Where(&Video{AuthorId: authorId}).Find(&data)
	//如果出现问题，返回对应到空，并且返回error
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

// Save 保存视频记录
func Save(videoName string, imageName string, authorId int64, title string) error {
	var video Video
	video.PublishTime = time.Now()
	video.PlayUrl = videoName
	video.CoverUrl = imageName
	video.AuthorId = authorId
	video.Title = title
	result := mysql.DB.Save(&video)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetVideosByLastTime
// 依据一个时间，来获取这个时间之前的一些视频
func GetVideosByLastTime(lastTime time.Time) ([]Video, error) {
	videos := make([]Video, lib.VideoCount)
	result := mysql.DB.Where("publish_time<?", lastTime).Order("publish_time desc").Limit(lib.VideoCount).Find(&videos)
	if result.Error != nil {
		return videos, result.Error
	}
	return videos, nil
}

// GetVideoIdsByAuthorId
// 通过作者id来查询发布的视频id切片集合
func GetVideoIdsByAuthorId(authorId int64) ([]int64, error) {
	var id []int64
	//通过pluck来获得单独的切片
	result := mysql.DB.Model(&Video{}).Where("author_id", authorId).Pluck("id", &id)
	//如果出现问题，返回对应到空，并且返回error
	if result.Error != nil {
		return nil, result.Error
	}
	return id, nil
}

// GetVideoByVideoId
// 依据VideoId来获得视频信息
func GetVideoByVideoId(videoId int64) (Video, error) {
	var tableVideo Video
	tableVideo.Id = videoId
	//Init()
	result := mysql.DB.First(&tableVideo)
	if result.Error != nil {
		return tableVideo, result.Error
	}
	return tableVideo, nil

}
