package controller



// 
// 
//   这个文件是一个测试文件，只是测试前端消息模块能否正确输出，实现了朋友列表的测试功能
//   并把router中的路由函数改成这个文件中的
// 
// 
// 


import (
	"net/http"

	"github.com/gin-gonic/gin"
)
type FriendUser struct {
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
	Message string `json:"message,omitempty"`
	MsgType int64  `json:"msgType,omitempty"`
}

type MessageUserListResponse struct {
	Response
	FriendList []FriendUser `json:"user_list"`
}

// 好友列表  POST /douyin/relation/friend/list/
// 这是一个好友列表的测试接口
func FriendsList(c *gin.Context) {
	// userId := c.Query("user_id")
	// uId, _ := strconv.ParseInt(userId, 10, 64)

	c.JSON(http.StatusOK, MessageUserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		FriendList: []FriendUser{
			FriendUser{Id: 4,Name: "测试用户",FollowCount: 10,FollowerCount: 10,IsFollow: true,
		Message: "test",MsgType: 0,},
			FriendUser{Id: 5,Name: "测试用户2",FollowCount: 10,FollowerCount: 10,IsFollow: true,
		Message: "tesst2",MsgType: 1,},
		},
	})
}