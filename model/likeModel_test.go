package model

import (
	"testing"

	"github.com/RaymondCode/simple-demo/lib"
	"github.com/RaymondCode/simple-demo/model/mysql"
)

func TestXxx(t *testing.T) {
	serverConfig := lib.LoadServerConfig()
	// 初始化数据库
	mysql.InitDB(serverConfig)

	InsertFavourite(Like{UserId: 1,VideoId: 2})
}