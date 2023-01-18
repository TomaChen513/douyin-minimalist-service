package model

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

func TestGetUserById(t *testing.T) {
	user,err:=GetUserById(2)
	fmt.Printf("%v", user)
	fmt.Printf("%v", err)
	user,err=GetUserById(10)
	fmt.Printf("%v", user)
	fmt.Printf("%v", err)
	user,err=GetUserById(0)
	fmt.Printf("%v", user)
	fmt.Printf("%v", err)
}


func TestGetUserByName(t *testing.T) {
	user,err:=GetUserByName("toma")
	fmt.Printf("%v", user)
	fmt.Printf("%v", err)
	user,err=GetUserByName("xxx")
	fmt.Printf("%v", user)
	fmt.Printf("%v", err)
	user,err=GetUserByName("")
	fmt.Printf("%v", user)
	fmt.Printf("%v", err)
}

func TestInsertUser(t *testing.T) {
	flag:=InsertUser(User{Name: "tommm",Password: "chenxinyu"})
	fmt.Printf("%v", flag)
	flag=InsertUser(User{Name: "tommm",Password: "chenxinyu"})
	fmt.Printf("%v", flag)
}