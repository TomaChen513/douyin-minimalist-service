package service

import (
	"strconv"
	"time"

	"github.com/RaymondCode/simple-demo/lib"
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
func (*FollowServiceImp) FollowAction(userId, followerId int64, cancel int8) bool {
	id := model.GetFollow(userId, followerId)
	if cancel == 1 {
		updateRedisWithAdd(followerId, userId)
		return model.UpdateFollow(id, cancel)
	} else {
		updateRedisWithDel(followerId, userId)
		return model.InsertFollow(userId, followerId, cancel)
	}

	//查询是否曾经关注, id为-1表示没有, 有则更新, 没有则插入一条新数据
	// if id := model.GetFollow(userId, followerId); id != -1 {
	// 	updateRedisWithAdd(followerId, userId)
	// 	return model.UpdateFollow(id, cancel)
	// } else if cancel == 1 {
	// 	return model.InsertFollow(userId, followerId, cancel)
	// } else {
	// 	return false
	// }
}

// 关注时, 更新redis
func updateRedisWithAdd(userId int64, targetId int64) (bool, error) {
	// 更新lib信息。
	/*
		1-Redis是否存在followers_targetId.
		2-Redis是否存在following_userId.
		3-Redis是否存在following_part_userId.
	*/
	// step1
	targetIdStr := strconv.Itoa(int(targetId))
	if cnt, _ := lib.RdbFollowers.SCard(lib.Ctx, targetIdStr).Result(); cnt != 0 {
		lib.RdbFollowers.SAdd(lib.Ctx, targetIdStr, userId)
		lib.RdbFollowers.Expire(lib.Ctx, targetIdStr, time.Hour*24)
	}
	// step2
	followingUserIdStr := strconv.Itoa(int(userId))
	if cnt, _ := lib.RdbFollowing.SCard(lib.Ctx, followingUserIdStr).Result(); cnt != 0 {
		lib.RdbFollowing.SAdd(lib.Ctx, followingUserIdStr, targetId)
		lib.RdbFollowing.Expire(lib.Ctx, followingUserIdStr, time.Hour*24)
	}
	// step3
	followingPartUserIdStr := followingUserIdStr
	lib.RdbFollowingPart.SAdd(lib.Ctx, followingPartUserIdStr, targetId)
	// 可能是第一次给改用户加followingPart的关注者，需要加上-1防止脏读。
	lib.RdbFollowingPart.SAdd(lib.Ctx, followingPartUserIdStr, -1)
	lib.RdbFollowingPart.Expire(lib.Ctx, followingPartUserIdStr, time.Hour*24)
	return true, nil
}

// 当取关时，更新lib里的信息
func updateRedisWithDel(userId int64, targetId int64) (bool, error) {
	/*
		1-Redis是否存在followers_targetId.
		2-Redis是否存在following_userId.
		2-Redis是否存在following_part_userId.
	*/
	// step1
	targetIdStr := strconv.Itoa(int(targetId))
	if cnt, _ := lib.RdbFollowers.SCard(lib.Ctx, targetIdStr).Result(); cnt != 0 {
		lib.RdbFollowers.SRem(lib.Ctx, targetIdStr, userId)
		lib.RdbFollowers.Expire(lib.Ctx, targetIdStr, time.Hour*24)
	}
	// step2
	followingIdStr := strconv.Itoa(int(userId))
	if cnt, _ := lib.RdbFollowing.SCard(lib.Ctx, followingIdStr).Result(); cnt != 0 {
		lib.RdbFollowing.SRem(lib.Ctx, followingIdStr, targetId)
		lib.RdbFollowing.Expire(lib.Ctx, followingIdStr, time.Hour*24)
	}
	// step3
	followingPartUserIdStr := followingIdStr
	if cnt, _ := lib.RdbFollowingPart.Exists(lib.Ctx, followingPartUserIdStr).Result(); cnt != 0 {
		lib.RdbFollowingPart.SRem(lib.Ctx, followingPartUserIdStr, targetId)
		lib.RdbFollowingPart.Expire(lib.Ctx, followingPartUserIdStr, time.Hour*24)
	}
	return true, nil
}

// 获取关注列表, userId表示查询对象, curId表示当前登录Id
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
	// 查看Redis中是否有关注数。
	if cnt, _ := lib.RdbFollowing.SCard(lib.Ctx, strconv.Itoa(int(userId))).Result(); cnt > 0 {
		// 更新过期时间。
		lib.RdbFollowing.Expire(lib.Ctx, strconv.Itoa(int(userId)), lib.ExpireTime)
		return cnt - 1, false
	}

	// 用SQL查询。
	ids, ok := model.GetFollowIds(userId)

	if !ok {
		return 0, false
	}
	// 更新Redis中的followers和followPart
	go addFollowingToRedis(int(userId), ids)

	return int64(len(ids)), true
}

func addFollowingToRedis(userId int, ids []int64) {
	lib.RdbFollowers.SAdd(lib.Ctx, strconv.Itoa(userId), -1)
	for i, id := range ids {
		lib.RdbFollowers.SAdd(lib.Ctx, strconv.Itoa(userId), id)
		lib.RdbFollowingPart.SAdd(lib.Ctx, strconv.Itoa(int(id)), userId)
		lib.RdbFollowingPart.SAdd(lib.Ctx, strconv.Itoa(int(id)), -1)
		// 更新部分关注者的时间
		lib.RdbFollowingPart.Expire(lib.Ctx, strconv.Itoa(int(id)),
			lib.ExpireTime+time.Duration((i%10)<<8))
	}
	// 更新followers的过期时间。
	lib.RdbFollowers.Expire(lib.Ctx, strconv.Itoa(userId), lib.ExpireTime)
}

// GetFollowerCnt 给定当前用户id，查询其粉丝数量。
func (*FollowServiceImp) GetFollowerCnt(userId int64) int64 {
	// 查Redis中是否已经存在。
	if cnt, _ := lib.RdbFollowers.SCard(lib.Ctx, strconv.Itoa(int(userId))).Result(); cnt > 0 {
		// 更新过期时间。
		lib.RdbFollowers.Expire(lib.Ctx, strconv.Itoa(int(userId)), time.Hour*24)
		return cnt - 1
	}
	//SQL中查询。
	ids, ok := model.GetFollowerIds(userId)
	if !ok {
		return 0
	}
	// 将数据存入Redis.
	// 更新followers 和 followingPart
	go addFollowersToRedis(int(userId), ids)

	return int64(len(ids))
}

func addFollowersToRedis(userId int, ids []int64) {
	lib.RdbFollowers.SAdd(lib.Ctx, strconv.Itoa(userId), -1)
	for i, id := range ids {
		lib.RdbFollowers.SAdd(lib.Ctx, strconv.Itoa(userId), id)
		lib.RdbFollowingPart.SAdd(lib.Ctx, strconv.Itoa(int(id)), userId)
		lib.RdbFollowingPart.SAdd(lib.Ctx, strconv.Itoa(int(id)), -1)
		// 更新部分关注者的时间
		lib.RdbFollowingPart.Expire(lib.Ctx, strconv.Itoa(int(id)),
			time.Hour*24+time.Duration((i%10)<<8))
	}
	// 更新followers的过期时间。
	lib.RdbFollowers.Expire(lib.Ctx, strconv.Itoa(userId), time.Hour*24)

}

// IsFollowing 给定当前用户和目标用户id，判断是否存在关注关系。
func (*FollowServiceImp) IsFollowing(userId int64, targetId int64) bool {
	// 先查Redis里面是否有此关系。
	if flag, _ := lib.RdbFollowingPart.SIsMember(lib.Ctx, strconv.Itoa(int(userId)), targetId).Result(); flag {
		// 重现设置过期时间。
		lib.RdbFollowingPart.Expire(lib.Ctx, strconv.Itoa(int(userId)), time.Hour*24)
		return true
	}
	//SQL 查询。
	if ok := model.IsFollow(targetId, userId); ok {
		// 存在此关系，将其注入Redis中。
		go addRelationToRedis(int(userId), int(targetId))
	}

	return false
}

func addRelationToRedis(userId int, targetId int) {
	// 第一次存入时，给该key添加一个-1为key，防止脏数据的写入。当然set可以去重，直接加，便于CPU。
	lib.RdbFollowingPart.SAdd(lib.Ctx, strconv.Itoa(int(userId)), -1)
	// 将查询到的关注关系注入Redis.
	lib.RdbFollowingPart.SAdd(lib.Ctx, strconv.Itoa(int(userId)), targetId)
	// 更新过期时间。
	lib.RdbFollowingPart.Expire(lib.Ctx, strconv.Itoa(int(userId)), time.Hour*24)
}
