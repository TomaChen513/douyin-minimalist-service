package lib

import (
	"github.com/go-redis/redis/v8"
)

// var Ctx = context.Background()
var RdbFollowers *redis.Client
var RdbFollowing *redis.Client
var RdbFollowingPart *redis.Client
