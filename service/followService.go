package service

import (
	"github.com/RaymondCode/simple-demo/model"
)

// 关注/取消关注操作
func FollowAction(userId, followId int64, cancel int8) bool {
	//查询是否曾经关注, id为-1表示没有, 有则更新, 没有则插入一条新数据
	if id := model.GetFollow(userId, followId); id != -1 {
		return model.UpdateFollow(id, cancel)
	} else if cancel == 1 {
		return model.InsertFollow(userId, followId, cancel)
	} else {
		return false
	}
}

// 获取关注列表, 失败返回false, userId表示查询对象, curId表示当前登录Id
func GetFollowList(userId, curId int64) ([]User, bool) {
	ids, ok := model.GetFollowIds(userId)
	if !ok {
		return nil, false
	}
	return GetUsersByids(ids, curId)
}

// 获取粉丝列表, 失败返回false, userId表示查询对象, curId表示当前登录Id
func GetFollowerList(userId, curId int64) ([]User, bool) {
	ids, ok := model.GetFollowerIds(userId)
	if !ok {
		return nil, false
	}

	return GetUsersByids(ids, curId)

}
