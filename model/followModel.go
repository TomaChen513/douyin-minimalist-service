package model

import (
	"log"

	"github.com/RaymondCode/simple-demo/model/mysql"
)

// Follow 用户关系结构，对应用户关系表。
type Follow struct {
	Id         int64 `gorm:"primaryKey"`
	UserId     int64
	FollowerId int64
	//1表示关注， 2表示不关注 user-->follower
	Cancel int8
}

// TableName 设置Follow结构体对应数据库表名。
func (Follow) TableName() string {
	return "follows"
}

// 查询user是否关注follower
func IsFollow(userId, followerId int64) bool {
	var cnt int64

	if err := mysql.DB.
		Model(Follow{}).
		Where("user_id = ?", userId).
		Where("follower_id = ?", followerId).
		Where("cancel = ?", 1).
		Count(&cnt).Error; err != nil {
		log.Println(err.Error())
		return false

	} else {
		return cnt != 0
	}
}

// 查询是否关注过
func GetFollow(userId, followerId int64) int64 {
	follow := Follow{}
	//存在返回Id
	err := mysql.DB.
		Where("user_id = ?", userId).
		Where("follower_id = ?", followerId).
		First(&follow).Error
	if err != nil {
		return -1
	}
	return follow.Id
}

// 插入一条新数据, 成功返回true， 失败返回false
func InsertFollow(userId, followerId int64, cancel int8) bool {
	follow := Follow{UserId: userId, FollowerId: followerId, Cancel: cancel}
	if err := mysql.DB.Create(&follow).Error; err != nil {
		log.Println(err.Error())
		return false
	}

	return true
}

// 修改数据, 成功返回true， 失败返回false
func UpdateFollow(id int64, cancel int8) bool {
	if err := mysql.DB.
		Model(&Follow{}).
		Where("id = ?", id).
		Update("cancel", cancel).Error; err != nil {
		log.Println(err.Error())
		return false
	}

	return true
}

// 查询粉丝人数
func GetFollowerCount(userId int64) int64 {
	var cnt int64
	//查询失败， 返回-1
	if err := mysql.DB.Model(Follow{}).
		Where("follower_id = ?", userId).
		Where("cancel = ?", 1).
		Count(&cnt).Error; err != nil {
		log.Println(err.Error())
		return -1
	}
	//查询成功， 返回粉丝数量
	return cnt
}

// 查询粉丝列表
func GetFollowerIds(userId int64) ([]int64, bool) {
	var ids []int64
	if err := mysql.DB.
		Model(Follow{}).
		Where("follower_id = ?", userId).
		Where("cancel = ?", 1).
		Pluck("user_id", &ids).Error; nil != err {
		// 没有关注任何人，但是不能算错。
		if err.Error() == "record not found" {
			return nil, true
		}
		// 查询失败
		log.Println(err.Error())
		return nil, false
	}
	// 查询成功
	return ids, true
}

// 查询关注数量
func GetFollowCount(userId int64) int64 {
	var cnt int64
	// 查询失败, 返回-1
	if err := mysql.DB.
		Model(Follow{}).
		Where("user_id = ?", userId).
		Where("cancel = ?", 1).
		Count(&cnt).Error; err != nil {
		log.Println(err.Error())
		return -1
	}
	// 查询成功，返回关注数量
	return cnt
}

// 查询关注列表
func GetFollowIds(userId int64) ([]int64, bool) {
	var ids []int64
	if err := mysql.DB.
		Model(Follow{}).
		Where("user_id = ?", userId).
		Where("cancel = ?", 1).
		Pluck("follower_id", &ids).Error; err != nil {
		// 没有粉丝，但是不能算错
		if err.Error() == "record not found" {
			return nil, true
		}
		// 查询出错
		log.Println(err.Error())
		return nil, false
	}
	// 查询成功
	return ids, true
}
