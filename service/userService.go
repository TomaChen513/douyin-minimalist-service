package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/RaymondCode/simple-demo/model"
)

type User struct {
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}

// 用户服务接口
type UserService interface {
	// 根据id获取用户
	GetUserById(id int64) (User, error)
	// 验证用户密码
	ValidPassword(id int64, password string) bool
	// 根据姓名获取用户
	GetUserByName(userName string) (User, error)
}

type UserServiceImpl struct {
	// FollowService
	// LikeService
}

// 根据id获取用户
func (usi *UserServiceImpl) GetUserById(id int64) (User, error) {
	tableUser, err := model.GetUserById(id)

	if err != nil {
		log.Println(err.Error())
		return User{}, err
	}

	return User{
		Id:            tableUser.Id,
		Name:          tableUser.Name,
		FollowCount:   100,
		FollowerCount: 100,
		IsFollow:      true,
	}, nil
}

// 验证用户密码
func (usi *UserServiceImpl) ValidPassword(id int64, password string) bool {
	tableUser, _ := model.GetUserById(id)
	return tableUser.Password == password
}

// 根据姓名获取用户
func (usi *UserServiceImpl) GetUserByName(userName string) (User, error) {
	tableUser, err := model.GetUserByName(userName)

	if err != nil {
		log.Println(err.Error())
		return User{}, err
	}
	return User{
		Id:            tableUser.Id,
		Name:          tableUser.Name,
		FollowCount:   0,
		FollowerCount: 0,
		IsFollow:      false,
	}, nil
}

// EnCoder 密码加密
func EnCoder(password string) string {
	h := hmac.New(sha256.New, []byte(password))
	sha := hex.EncodeToString(h.Sum(nil))
	fmt.Println("Result: " + sha)
	return sha
}

// InsertTableUser 将tableUser插入表内
func (usi *UserServiceImpl) InsertTableUser(User *model.User) bool {
	flag := model.InsertTableUser(User)
	if !flag {
		log.Println("插入失败")
		return false
	}
	return true
}
