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

	// 上传本地文件。
	err = bucket.PutObject(videoName, file)
	if err != nil {
		fmt.Println("本地文件上传Error:", err)
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
