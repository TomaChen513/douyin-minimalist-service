package service

import (
	"log"

	"github.com/RaymondCode/simple-demo/model"
)

type UserService interface {
	GetUserById(id int64) (User, error)
	GetUserByName(userName string) (User, error)
	InsertUser(user User, password string) bool
	ValidPassword(id int64, password string) bool
}

type User struct {
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}

type UserServiceImpl struct {
}

func (usi *UserServiceImpl) GetUserById(id int64) (User, error) {
	tableUser, err := model.GetUserById(id)
	var user User
	if err != nil {
		log.Println(err.Error())
		return user, err
	}
	return parseToUser(tableUser), nil
}

func (usi *UserServiceImpl) GetUserByName(userName string) (User, error) {
	tableUser, err := model.GetUserByName(userName)
	var user User
	if err != nil {
		log.Println(err.Error())
		return user, err
	}
	return parseToUser(tableUser), nil
}

func (usi *UserServiceImpl) InsertUser(user User, password string) bool {
	tableUser := model.TableUser{
		Id:            user.Id,
		Name:          user.Name,
		Password:      password,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		IsFollow:      user.IsFollow,
	}
	return model.InsertUser(tableUser)
}

func (usi *UserServiceImpl) ValidPassword(id int64, password string) bool {
	tableUser, _ := model.GetUserById(id)
	return tableUser.Password == password
}

// 将Table结构转换成json的结构
func parseToUser(tableUser model.TableUser) User {
	return User{
		Id:            tableUser.Id,
		Name:          tableUser.Name,
		FollowCount:   tableUser.FollowCount,
		FollowerCount: tableUser.FollowerCount,
		IsFollow:      tableUser.IsFollow,
	}
}
