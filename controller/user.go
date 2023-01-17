package controller

import (
	"net/http"
	// "sync/atomic"

	"github.com/RaymondCode/simple-demo/model"
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


// 取消了redis存储，采用map直接存储，在服务重启后清空
var userToken =map[string]int64{
	"zhangleidouyin":1,
}

// var userIdSequence = int64(1)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User model.User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	token := username + password

	// id, _ := lib.GetKey(token)
	// 用map来提取id
	id:=userToken[token]

	if _, exist := model.QueryUserExists(id); exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		// 创建新用户
		newUser:=model.CreateUser(username,password)
		// //存入redis
		// if err := lib.SetKey(token, strconv.Itoa(int(newUser.Id)), 3600); err != nil {
		// 	return
		// }
		// 将redis改为map
		userToken[token]=newUser.Id

		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   newUser.Id,
			Token:    username + password,
		})
	}

	// if _, exist := usersLoginInfo[token]; exist {
	// 	c.JSON(http.StatusOK, UserLoginResponse{
	// 		Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
	// 	})
	// } else {
	// 	atomic.AddInt64(&userIdSequence, 1)
	// 	newUser := User{
	// 		Id:   userIdSequence,
	// 		Name: username,
	// 	}
	// 	usersLoginInfo[token] = newUser
	// 	c.JSON(http.StatusOK, UserLoginResponse{
	// 		Response: Response{StatusCode: 0},
	// 		UserId:   userIdSequence,
	// 		Token:    username + password,
	// 	})
	// }
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	token := username + password
	// token:="zhangleidouyin"

	user := model.VerifyPasswd(username, password)


	if user == nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	} else {
		//存入redis
		// if err := lib.SetKey(token, strconv.Itoa(int(user.Id)), 3600); err != nil {
		// 	return
		// }
		userToken[token]=user.Id

		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    token,
		})
	}

	// if user, exist := usersLoginInfo[token]; exist {
	// 	c.JSON(http.StatusOK, UserLoginResponse{
	// 		Response: Response{StatusCode: 0},
	// 		UserId:   user.Id,
	// 		Token:    token,
	// 	})
	// } else {
	// 	c.JSON(http.StatusOK, UserLoginResponse{
	// 		Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
	// 	})
	// }
}

func UserInfo(c *gin.Context) {
	token := c.Query("token")

	// id, _ := lib.GetKey(token)
	id:=userToken[token]

	if user, exist := model.QueryUserExists(id); exist {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     *user,
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
	// if user, exist := usersLoginInfo[token]; exist {
	// 	c.JSON(http.StatusOK, UserResponse{
	// 		Response: Response{StatusCode: 0},
	// 		User:     user,
	// 	})
	// } else {
	// 	c.JSON(http.StatusOK, UserResponse{
	// 		Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
	// 	})
	// }
}
