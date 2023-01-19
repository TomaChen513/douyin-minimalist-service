package service

import (
	"fmt"
	"os"
	"testing"

	"github.com/RaymondCode/simple-demo/lib"
	"github.com/RaymondCode/simple-demo/model/mysql"
)

func TestMain(m *testing.M){
	serverConfig := lib.LoadServerConfig()
	mysql.InitDB(serverConfig)
	code:=m.Run()
	os.Exit(code)
}

func TestGetVideosByUserId(t *testing.T) {
	vsi:=VideoServiceImpl{}
	user,err:=vsi.GetVideosByUser(2)
	fmt.Printf("%+v  ", user)
	fmt.Printf("%v\n", err)
	user,err=vsi.GetVideosByUser(3)
	fmt.Printf("%+v  ", user)
	fmt.Printf("%v\n", err)
	user,err=vsi.GetVideosByUser(-1)
	fmt.Printf("%+v  ", user)
	fmt.Printf("%v\n", err)
}