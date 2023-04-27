package lib

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

func TestSetKey(t *testing.T) {
	var Ctx = context.Background()

	RdbLikeUserId = redis.NewClient(&redis.Options{
		Addr:     "121.5.231.228:6379",
		Password: "wintercamp",
		DB:       0, //  选择将点赞视频id信息存入 DB0.
	})

	RdbLikeUserId.Set(Ctx,"toma","cxy1",20*time.Second)


	// RdbLikeUserId.SAdd(Ctx,"xyz","123","456","789")

	
}

func TestGetKey(t *testing.T) {
	var Ctx = context.Background()

	RdbLikeUserId = redis.NewClient(&redis.Options{
		Addr:     "121.5.231.228:6379",
		Password: "wintercamp",
		DB:       0, //  选择将点赞视频id信息存入 DB0.
	})

	RdbLikeUserId.Set(Ctx,"toma","cxy1",20*time.Second)


	// RdbLikeUserId.SAdd(Ctx,"xyz","123","456","789")

	
}
