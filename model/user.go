package model

import (
	"github.com/RaymondCode/simple-demo/model/mysql"
	"github.com/jinzhu/gorm"
)

type User struct {
	Id            int64  `gorm:"primarykey" json:"id,omitempty"`
	Name          string `gorm:"type:varchar(255)" json:"name,omitempty"`
	Password      string `gorm:"type:varchar(20)" json:"password,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}

// 查询判断用户是否存在
func QueryUserExists(id string) bool {
	var user User
	mysql.DB.Find(&user, "id = ?", id)
	if user.Id == 0 {
		return false
	}
	return true
}

func VerifyPasswd(userName, passWord string) int64 {
	var user []User
	err := mysql.DB.Model(&User{}).
		Where(&User{Name: userName, Password: passWord}).
		Find(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0
	}
	return user[0].Id
}
