package service

import (
	"log"

	"github.com/RaymondCode/simple-demo/model"
)

type User struct {
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
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
	FollowService
	FavorService
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
		FollowCount:   0,
		FollowerCount: 0,
		IsFollow:      false,
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

// 根据id获得用户详细信息, curId表示当前登录的Id
func (usi *UserServiceImpl) GetUserInfoById(id, curId int64) (User, bool) {
	name := model.GetNameById(id)
	//println(name)
	if name == "" {
		return User{}, false
	}

	followCount := model.GetFollowCount(id)
	//println(followCount)
	if followCount == -1 {
		return User{}, false
	}

	followerCount := model.GetFollowerCount(id)
	//println(followerCount)
	if followerCount == -1 {
		return User{}, false
	}

	isFollow := model.IsFollow(curId, id)
	println(isFollow)
	if !isFollow {
		return User{}, false
	}

	return User{
		Id:            id,
		Name:          name,
		FollowCount:   followCount,
		FollowerCount: followerCount,
		IsFollow:      isFollow,
	}, true
	//return User{}, true
}

// 根据ids获取users, curId表示当前用户id
func (usi *UserServiceImpl) GetUsersByids(ids []int64, curId int64) ([]User, bool) {
	var users = make([]User, 0, len(ids))

	for _, id := range ids {
		if user, ok := usi.GetUserInfoById(id, curId); !ok {
			return nil, false
		} else {
			users = append(users, user)
		}
		//println(id)
	}

	return users, true
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
