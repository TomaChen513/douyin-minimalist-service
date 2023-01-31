package main

import (
	"log"

	"github.com/RaymondCode/simple-demo/lib"
	"github.com/RaymondCode/simple-demo/model/mysql"
	"github.com/RaymondCode/simple-demo/router"
)

func main() {
	// 加载配置
	serverConfig := lib.LoadServerConfig()
	// 初始化数据库
	mysql.InitDB(serverConfig)
	defer mysql.DB.Close()

	// 设置路由
	r := router.SetupRoute()

	// 启动服务
	if err := r.Run(); err != nil {
		log.Fatal("服务器启动失败...")
	}

}
