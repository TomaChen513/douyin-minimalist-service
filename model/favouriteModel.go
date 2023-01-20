package model

import (
	"log"

	"github.com/RaymondCode/simple-demo/model/mysql"
)

type TableFavourite struct {
	Id int64 
	UserId int64
	VideoId int64
}

// TableName 修改表名映射
func (tableFavourite TableFavourite) TableName() string {
	return "favourite"
}

// 点赞，执行后数据库添加一条点赞记录
func InsertFavourite(favourite TableFavourite) bool{
	// 可以重复点赞，需要修复
	if err:=mysql.DB.Create(&favourite).Error;err!=nil {
		log.Println(err.Error())
		return false
	}
	return true
}

// 取消点赞，执行后删除点赞记录
func DeleteFavourite(favourite TableFavourite) bool{
	if err:=mysql.DB.Where("user_id= ?",favourite.UserId).Where("video_id=?",favourite.VideoId).Delete(&favourite).Error;err!=nil {
		log.Println(err.Error())
		return false
	}
	return true
}

// 根据用户id寻找视频
func SelectVideosByUserId(userId int64) ([]int64,error){
	var favouriteVideos []TableFavourite
	if err:=mysql.DB.Where("user_id=?",userId).Find(&favouriteVideos).Error;err!=nil {
		log.Println(err.Error())
		return []int64{},err
	}

	res:=make([]int64,len(favouriteVideos))
	for i := 0; i < len(favouriteVideos); i++ {
		res[i]=favouriteVideos[i].VideoId
	}
	return res,nil
}