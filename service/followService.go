package service

import (
	"github.com/RaymondCode/simple-demo/model"
)

type FollowService interface {
	FollowAction(userId, followId int64, cancel int8) bool
	GetFollowList(userId, curId int64) ([]User, bool)
	GetFollowerList(userId, curId int64) ([]User, bool)
	// GetFollowingCnt 根据用户id来查询用户关注了多少其它用户
	GetFollowingCnt(userId int64) (int64, bool)
	// GetFollowerCnt 根据用户id来查询用户被多少其他用户关注
	GetFollowerCnt(userId int64) int64
	// IsFollowing 根据当前用户id和目标用户id来判断当前用户是否关注了目标用户
	IsFollowing(userId int64, targetId int64) bool
}

type FollowServiceImp struct {
	UserService
}

// 关注/取消关注操作
func (fsi *FollowServiceImp) FollowAction(userId, followId int64, cancel int8) bool {
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
func (fsi *FollowServiceImp) GetFollowList(userId, curId int64) ([]User, bool) {
	ids, ok := model.GetFollowIds(userId)
	if !ok {
		return nil, false
	}

	return fsi.GetUsersByids(ids, curId)
}

// 获取粉丝列表, 失败返回false, userId表示查询对象, curId表示当前登录Id
func (fsi *FollowServiceImp) GetFollowerList(userId, curId int64) ([]User, bool) {
	ids, ok := model.GetFollowerIds(userId)
	if !ok {
		return nil, false
	}

	return fsi.GetUsersByids(ids, curId)

}

func (fsi *FollowServiceImp) GetFriendList(userId, curId int64) ([]User, bool) {
	ids, ok := model.GetFriendIds(userId)
	if !ok {
		return nil, false
	}

	return fsi.GetUsersByids(ids, curId)
}

// GetFollowingCnt 给定当前用户id，查询其关注者数量。
func (*FollowServiceImp) GetFollowingCnt(userId int64) (int64, bool) {
	// // 查看Redis中是否有关注数。
	// if cnt, err := redis.RdbFollowing.SCard(redis.Ctx, strconv.Itoa(int(userId))).Result(); cnt > 0 {
	// 	// 更新过期时间。
	// 	redis.RdbFollowing.Expire(redis.Ctx, strconv.Itoa(int(userId)), config.ExpireTime)
	// 	return cnt - 1, err
	// }

	// 用SQL查询。
	ids, ok := model.GetFollowIds(userId)

	if !ok {
		return 0, false
	}
	// // 更新Redis中的followers和followPart
	// go addFollowingToRedis(int(userId), ids)

	return int64(len(ids)), true
}

// GetFollowerCnt 给定当前用户id，查询其粉丝数量。
func (*FollowServiceImp) GetFollowerCnt(userId int64) int64 {
	// // 查Redis中是否已经存在。
	// if cnt, err := redis.RdbFollowers.SCard(redis.Ctx, strconv.Itoa(int(userId))).Result(); cnt > 0 {
	// 	// 更新过期时间。
	// 	redis.RdbFollowers.Expire(redis.Ctx, strconv.Itoa(int(userId)), config.ExpireTime)
	// 	return cnt - 1, err
	// }
	// SQL中查询。
	cnt := model.GetFollowerCount(userId)
	// if nil != err {
	// 	return 0, err
	// }
	// // 将数据存入Redis.
	// // 更新followers 和 followingPart
	// go addFollowersToRedis(int(userId), ids)

	return cnt
}

// IsFollowing 给定当前用户和目标用户id，判断是否存在关注关系。
func (*FollowServiceImp) IsFollowing(userId int64, targetId int64) bool {
	// // 先查Redis里面是否有此关系。
	// if flag, err := redis.RdbFollowingPart.SIsMember(redis.Ctx, strconv.Itoa(int(userId)), targetId).Result(); flag {
	// 	// 重现设置过期时间。
	// 	redis.RdbFollowingPart.Expire(redis.Ctx, strconv.Itoa(int(userId)), config.ExpireTime)
	// 	return true, err
	// }
	// SQL 查询。
	cnt := model.IsFollow(userId, targetId)

	// // 存在此关系，将其注入Redis中。
	// go addRelationToRedis(int(userId), int(targetId))

	return cnt
}
