package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/RaymondCode/simple-demo/lib"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User service.User `json:"user"`
}

// Register POST douyin/user/register/ 用户注册
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	usi:=service.UserServiceImpl{}

	// 根据注册的用户名查找是否有相同用户
	user, _ := usi.GetUserByName(username)

	// 重名情况
	if user.Name == username {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户已经存在"},
		})
	} else {
		// 创建新用户
		newUser := service.User{Name: username}
		if !usi.InsertUser(newUser,password) {
			log.Println("创建新用户失败")
			return
		}
		newUser,_ = usi.GetUserByName(newUser.Name)
		// 产生新token
		token,_ := service.ReleaseToken(newUser)
		// 将token存入redis
		if err := lib.SetKey(token, strconv.Itoa(int(newUser.Id)), 3600); err != nil {
			log.Println("token存入redis失败")
			return
		}
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   newUser.Id,
			Token:    token,
		})
	}
}

// Login POST douyin/user/login/ 用户登录
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	usi:=service.UserServiceImpl{}

	user, err := usi.GetUserByName(username)

	
	// 若用户不存在
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户不存在"},
		})
		return
	}

	// 验证密码
	if  usi.ValidPassword(user.Id,password){
		// 分发token
		token,_ := service.ReleaseToken(user)
		// 存入redis
		if err := lib.SetKey(token, strconv.Itoa(int(user.Id)), 3600); err != nil {
			log.Println("token存入redis失败")
			return
		}
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    token,
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "密码错误"},
		})
	}
}

// 获得用户信息 GET /douyin/user/
func UserInfo(c *gin.Context) {
	token := c.Query("token")
	userId:=c.Query("user_id")
	id,_:=strconv.ParseInt(userId,10,64)

	usi:=service.UserServiceImpl{}

	// 验证token
	tId,_:=lib.GetKey(token)
	tokenId,_:=strconv.ParseInt(tId,10,64)
	if id!= tokenId{
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户id与token信息不一致！"},
		})
		return
	}

	if user,err:=usi.GetUserById(id);err!=nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户不存在"},
		})
	}else{
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     user,
		})
	}
}
