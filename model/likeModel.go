package model

import (
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