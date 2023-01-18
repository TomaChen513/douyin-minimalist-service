package model

import (
	"errors"
	"log"

	"github.com/RaymondCode/simple-demo/model/mysql"
)

type User struct {
	Id            int64  `gorm:"primarykey" json:"id,omitempty"`
	Name          string `gorm:"type:varchar(255)" json:"name,omitempty"`
	Password      string `gorm:"type:varchar(20)" json:"password,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}

// 根据id查询用户是否存在
func GetUserById(id int64) (User, error) {
	user := User{}
	mysql.DB.First(&user, id)
	if user.Id == 0 {
		err := errors.New("找不到指定id的用户")
		log.Println(err.Error())
		return user, err
	}
	return user, nil
}

// 根据用户名获得user对象
func GetUserByName(userName string) (User, error) {
	user := User{}
	if err := mysql.DB.Where("name = ?", userName).Find(&user).Error; err != nil {
		log.Println(err.Error())
		return user, err
	}
	return user, nil
}

// 插入用户
func InsertUser(user User) bool {
	if err := mysql.DB.Create(&user).Error; err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}
