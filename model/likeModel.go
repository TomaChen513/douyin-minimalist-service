package model

import (
	"errors"
	"log"

	"github.com/RaymondCode/simple-demo/model/mysql"
)

type Like struct {
	Id      int64 `gorm:"primarykey" json:"id,omitempty"`
	UserId  int64 `json:"user_id,omitempty"`
	VideoId int64 `json:"video_id,omitempty"`
	// cancel字段暂时不使用
	Cancel  int64 
}


// TableName 修改表名映射
func (like Like) TableName() string {
	return "likes"
}

// 点赞，执行后数据库添加一条点赞记录
func InsertFavourite(like Like) bool{
	// 可以重复点赞，需要修复
	if err:=mysql.DB.Create(&like).Error;err!=nil {
		log.Println(err.Error())
		return false
	}
	return true
}

// 更新点赞记录
func UpdateLike(userId int64, videoId int64, actionType int32) error {
	//更新当前用户观看视频的点赞状态“cancel”，返回错误结果
	err := mysql.DB.Model(Like{}).Where(map[string]interface{}{"user_id": userId, "video_id": videoId}).
		Update("cancel", actionType).Error
	//如果出现错误，返回更新数据库失败
	if err != nil {
		log.Println(err.Error())
		return errors.New("update data fail")
	}
	//更新操作成功
	return nil
}

// 取消点赞，执行后删除点赞记录
func DeleteFavourite(like Like) bool{
	if err:=mysql.DB.Where("user_id= ?",like.UserId).Where("video_id=?",like.VideoId).Delete(&like).Error;err!=nil {
		log.Println(err.Error())
		return false
	}
	return true
}

// 根据用户id寻找视频标号
func SelectVideosByUserId(userId int64) ([]int64,error){
	var likes []Like
	if err:=mysql.DB.Where("user_id=?",userId).Find(&likes).Error;err!=nil {
		log.Println(err.Error())
		return []int64{},err
	}

	res:=make([]int64,len(likes))
	for i := 0; i < len(likes); i++ {
		res[i]=likes[i].VideoId
	}
	return res,nil
}

// 根据视频id寻找点赞记录
func CountLikesByVideoId(videoId int64) (int64,error){
	var likes []Like
	if err:=mysql.DB.Where("video_id=?",videoId).Find(&likes).Error;err!=nil {
		log.Println(err.Error())
		return -1,err
	}
	return int64(len(likes)),nil
}

// 根据userId和videoId查找
func GetLike(userId,videoId int64) (Like,error){
	like:=Like{}

	if err:=mysql.DB.Where("user_id = ? AND video_id = ?", userId, videoId).
	Find(&like).Error;err!=nil{
		log.Println(err.Error())
		return Like{},err
	}
	return like,nil
}

