package main

import (
	"log"

	"github.com/RaymondCode/simple-demo/lib"
	"github.com/RaymondCode/simple-demo/model/mysql"
	"github.com/RaymondCode/simple-demo/router"
	"github.com/RaymondCode/simple-demo/service"
)

func main() {
	go service.RunMessageServer()

	serverConfig := lib.LoadServerConfig()
	mysql.InitDB(serverConfig)
	defer mysql.DB.Close()

	r := router.SetupRoute()

	if err := r.Run(); err != nil {
		log.Fatal("服务器启动失败...")
	}
}
