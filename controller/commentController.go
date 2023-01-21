package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/RaymondCode/simple-demo/lib"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
)

type CommentListResponse struct {
	Response
	CommentList []service.Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment service.Comment `json:"comment,omitempty"`
}

// 评论操作  POST /douyin/comment/action/
func CommentAction(c *gin.Context) {
	token := c.Query("token")
	video_id := c.Query("video_id")
	actionType := c.Query("action_type")
	vId, err := strconv.ParseInt(video_id, 10, 64)

	if err != nil {
		log.Println("video_id读取失败")
		return
	}

	// 验证token
	tId, err := lib.GetKey(token)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "token信息有误！"},
		})
		return
	}
	uId, _ := strconv.ParseInt(tId, 10, 64)

	csi := service.CommentServiceImpl{}

	if actionType == "1" {
		// 发布评论
		content := c.Query("comment_text")
		comment, err := csi.InsertComment(uId, vId, content)
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "评论发布失败",
			})
			return
		}
		c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 0},
			Comment: comment})
	} else if actionType == "2" {
		comment_id := c.Query("comment_id")
		cId, _ := strconv.ParseInt(comment_id, 10, 64)
		// 删除评论
		if !csi.DeleteComment(cId) {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "评论删除失败"})
			return
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "actionType有误"})
	}
}

// 倒序获得所有评论  Get /douyin/comment/list/
func CommentList(c *gin.Context) {
	token := c.Query("token")
	video_id := c.Query("video_id")
	vId, _ := strconv.ParseInt(video_id, 10, 64)
	// 验证token
	_, err := lib.GetKey(token)
	if  err != nil{
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "token有误"})
		return
	}

	csi := service.CommentServiceImpl{}

	commentLists,err:=csi.GetCommentListDecByTime(vId)
	if err!=nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获得评论出错！"})
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: commentLists,
	})
}
