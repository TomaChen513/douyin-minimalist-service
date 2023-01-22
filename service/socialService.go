package service

import (
	"github.com/RaymondCode/simple-demo/model"
)

type SocialService interface {
	AttentionAction(userId, toUserId int64, actionType string) bool
	GetFollowList(userId int64) ([]User, error)
	GetFollowerList(userId int64) ([]User, error)
	GetFriendList(userId int64) ([]User, error)
}

type SocialServiceImpl struct {
	UserServiceImpl
}

// 重复问题没考虑
func (ssi *SocialServiceImpl) AttentionAction(userId, toUserId int64, actionType string) bool {
	tableRelation := model.TableRelation{UserId: userId, ToUserId: toUserId}
	if actionType == "1" {
		// 关注  改成事务操作
		ok := model.InsertRelation(tableRelation)
		if !ok {
			return false
		}
		return ssi.UpdateFollow(userId, toUserId, "1")
	}
	// 取消关注  改成事务操作
	ok := model.DeleteRelation(tableRelation)
	if !ok {
		return false
	}
	return ssi.UpdateFollow(userId, toUserId, "2")
}

func (ssi *SocialServiceImpl) GetFollowList(userId int64) ([]User, error) {
	tableRelations, err := model.GetRelationsByUserId(userId)
	followList := make([]User, len(tableRelations))
	if err != nil {
		return followList, err
	}
	for i := 0; i < len(tableRelations); i++ {
		user, _ := ssi.GetUserById(tableRelations[i].ToUserId)
		followList[i] = user
	}
	return followList, nil
}

func (ssi *SocialServiceImpl) GetFollowerList(toUserId int64) ([]User, error) {
	tableRelations, err := model.GetRelationsByToUserId(toUserId)
	followerList := make([]User, len(tableRelations))
	if err != nil {
		return followerList, err
	}
	for i := 0; i < len(tableRelations); i++ {
		user, _ := ssi.GetUserById(tableRelations[i].UserId)
		followerList[i] = user
	}
	return followerList, nil
}

// 正确好友数量显示bug
func (ssi *SocialServiceImpl) GetFriendList(userId int64) ([]User, error) {
	tableRelations, err := model.GetRelationsByUserId(userId)
	friendList := make([]User, 0)
	if err != nil {
		return friendList, err
	}
	for i := 0; i < len(tableRelations); i++ {
		// 查找是否有对应的关系
		ok := model.GetTableRelation(tableRelations[i].ToUserId, tableRelations[i].UserId)
		if ok {
			user, _ := ssi.GetUserById(tableRelations[i].ToUserId)
			friendList = append(friendList, user)
		}
	}
	return friendList, nil
}
