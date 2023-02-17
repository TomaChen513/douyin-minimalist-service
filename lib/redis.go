package lib

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

var RdbFollowers *redis.Client
var RdbFollowing *redis.Client
var RdbFollowingPart *redis.Client
var RdbLikeUserId *redis.Client     //key:userId,value:VideoId 用户点赞对应视频
var RdbLikeVideoCount *redis.Client //key:videoId,value:cnt 视频对应点赞数量
var RdbCommentVideoId *redis.Client //key:VideoId,value:userId 视频对应所评论的用户
var RdbCommentid *redis.Client      //key:videoId_userId,value:content
var RdbCommentCount *redis.Client   //key:videoId,value:cnt 评论对应评论数量

func init() {
	config := LoadServerConfig()

	RdbLikeUserId = redis.NewClient(&redis.Options{
		Addr:     config.RedisHost,
		Password: "wintercamp",
		DB:       0, //  选择将点赞视频id信息存入 DB0.
	})
	RdbCommentVideoId = redis.NewClient(&redis.Options{
		Addr:     config.RedisHost,
		Password: "wintercamp",
		DB:       1, //  选择将点赞用户id信息存入 DB1.
	})
	RdbCommentid = redis.NewClient(&redis.Options{
		Addr:     config.RedisHost,
		Password: "wintercamp",
		DB:       2, // 选择将video_id中的评论id s存入 DB2.
	})

	RdbCommentCount = redis.NewClient(&redis.Options{
		Addr:     config.RedisHost,
		Password: "wintercamp",
		DB:       3, // 选择将comment_id对应video_id存入 DB3.
	})

	RdbLikeVideoCount = redis.NewClient(&redis.Options{
		Addr:     config.RedisHost,
		Password: "wintercamp",
		DB:       4, // videoId对应的点赞数量存入DB4
	})
	RdbFollowers = redis.NewClient(&redis.Options{
		Addr:     config.RedisHost,
		Password: "wintercamp",
		DB:       5, // videoId对应的点赞数量存入DB4
	})
	RdbFollowing = redis.NewClient(&redis.Options{
		Addr:     config.RedisHost,
		Password: "wintercamp",
		DB:       6, // videoId对应的点赞数量存入DB4
	})
	RdbFollowingPart = redis.NewClient(&redis.Options{
		Addr:     config.RedisHost,
		Password: "wintercamp",
		DB:       7, // videoId对应的点赞数量存入DB4
	})
}
