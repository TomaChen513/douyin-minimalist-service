package lib

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

var RdbLikeUserId *redis.Client  //key:userId,value:VideoId 用户点赞对应视频
var RdbLikeVideoId *redis.Client //key:VideoId,value:userId 视频对应所点赞的用户
var RdbLikeVideoCount *redis.Client //key:videoId,value:cnt 视频对应点赞数量
var RdbVCid *redis.Client
var RdbCVid *redis.Client

func init() {
	config := LoadServerConfig()

	RdbLikeUserId = redis.NewClient(&redis.Options{
		Addr:     config.RedisHost,
		Password: "wintercamp",
		DB:       0, //  选择将点赞视频id信息存入 DB0.
	})
	RdbLikeVideoId = redis.NewClient(&redis.Options{
		Addr:     config.RedisHost,
		Password: "wintercamp",
		DB:       1, //  选择将点赞用户id信息存入 DB1.
	})
	RdbVCid = redis.NewClient(&redis.Options{
		Addr:     config.RedisHost,
		Password: "wintercamp",
		DB:       2, // 选择将video_id中的评论id s存入 DB2.
	})

	RdbCVid = redis.NewClient(&redis.Options{
		Addr:     config.RedisHost,
		Password: "wintercamp",
		DB:       3, // 选择将comment_id对应video_id存入 DB3.
	})

	RdbLikeVideoCount = redis.NewClient(&redis.Options{
		Addr:     config.RedisHost,
		Password: "wintercamp",
		DB:       4, // videoId对应的点赞数量存入DB4
	})
}
