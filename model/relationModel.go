package model

import (
	"log"

	"github.com/RaymondCode/simple-demo/model/mysql"
)

type TableRelation struct {
	Id, UserId, ToUserId int64
}

func (tableRelation TableRelation) TableName() string {
	return "relation"
}

func GetRelationsByUserId(userId int64) ([]TableRelation, error) {
	var tableRelations []TableRelation
	if err := mysql.DB.Where("user_id=?", userId).Find(&tableRelations).Error; err != nil {
		log.Println(err.Error())
		return tableRelations, err
	}
	return tableRelations, nil
}

func GetRelationsByToUserId(toUserId int64) ([]TableRelation, error) {
	var tableRelations []TableRelation
	if err := mysql.DB.Where("to_user_id=?", toUserId).Find(&tableRelations).Error; err != nil {
		log.Println(err.Error())
		return tableRelations, err
	}
	return tableRelations, nil
}

// 查找指定的关注者与被关注者
func GetTableRelation(userId, toUserId int64) bool {
	var tableRelation TableRelation
	if err := mysql.DB.Where("to_user_id=?", toUserId).
		Where("user_id=?", userId).
		Find(&tableRelation).Error; err!=nil {
		log.Println("无对应关系")
		return false
	}
	return true
}

func InsertRelation(tableRelation TableRelation) bool {
	if err := mysql.DB.Create(&tableRelation).Error; err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

func DeleteRelation(tableRelation TableRelation) bool {
	if err := mysql.DB.Where("user_id= ?", tableRelation.UserId).
		Where("to_user_id=?", tableRelation.ToUserId).
		Delete(&tableRelation).Error; err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}
