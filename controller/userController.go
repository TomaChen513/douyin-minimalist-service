package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User service.User `json:"user"`
}

// Login POST douyin/user/login/ 用户登录
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	// 新建用户实例
	usi := service.UserServiceImpl{}

	// 根据用户姓名获得密码
	user, err := usi.GetUserByName(username)

	// 若用户不存在
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户不存在"},
		})
		return
	}

	// 与传入的密码进行验证
	if usi.ValidPassword(user.Id, password) {
		// 验证成功则分发jwt token
		token, _ := service.ReleaseToken(user)

		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			// StatusMsg: "登录成功",
			UserId: user.Id,
			Token:  token,
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "密码错误"},
		})
	}
}

// Register POST douyin/user/register/ 用户注册
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	// 新建用户实例
	usi := service.UserServiceImpl{}

	// 根据用户姓名获得密码
	user, _ := usi.GetUserByName(username)

	if username == user.Name {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		newUser := model.User{
			Name: username,
			// Password: service.EnCoder(password),
			Password: password,
		}
		if !usi.InsertTableUser(&newUser) {
			println("Insert Data Fail")
		}
		// u := usi.GetTableUserByUsername(username)
		token, _ := service.ReleaseToken(user)
		log.Println("注册返回的id: ", user.Id)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    token,
		})
	}
}

// UserInfo GET douyin/user/ 用户信息
func UserInfo(c *gin.Context) {
	user_id := c.Query("user_id")
	user_token := c.Query("token")

	_, flag := service.ParseToken(user_token)
	if !flag {
		c.JSON(http.StatusUnauthorized, Response{
			StatusCode: -1,
			StatusMsg:  "Token Error",
		})
		return
	}

	id, _ := strconv.ParseInt(user_id, 10, 64)

	// 新建用户实例
	usi := service.UserServiceImpl{}

	// usi := service.UserServiceImpl{
	// 	FollowService: &service.FollowServiceImp{},
	// 	LikeService:   &service.LikeServiceImpl{},
	// }

	u, err := usi.GetUserById(id)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User Doesn't Exist"},
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     u,
		})
	}
}
