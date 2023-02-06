package controller

import (
	"log"
	"net/http"
	"strconv"

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
	user_id, exist := c.Get("userId")
	if !exist {
		log.Println("router:CommentAction jwt验证出错!")
		return
	}
	uId := user_id.(int64)
	vId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		log.Println("router:CommentAction video_id读取失败")
		return
	}
	actionType := c.Query("action_type")

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
		c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 0,StatusMsg: "评论发布成功!"},
			Comment: comment})
	} else if actionType == "2" {
		// 删除评论
		cId, err := strconv.ParseInt(c.Query("comment_id"), 10, 64)
		if err != nil {
			log.Println("router:CommentAction comment_id解析失误!")
		}

		if !csi.DeleteComment(cId,uId,vId) {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "评论删除失败"})
			return
		}
		// 成功返回
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		log.Println("router:CommentAction actionType错误!")
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "actionType有误"})
	}
}

// 倒序获得所有评论  Get /douyin/comment/list/
func CommentList(c *gin.Context) {
	vId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "video_id格式有误",
		})
		log.Println("router:CommentList video_id格式有误") //视频id格式有误
		return
	}

	csi := service.CommentServiceImpl{}

	commentLists, err := csi.GetCommentListDecByTime(vId)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获得评论出错！"})
		log.Println("router:csi.GetCommentListDecByTime(vId) 服务调用出错") //视频id格式有误
		return
	}

	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0, StatusMsg: "成功获得视频评论"},
		CommentList: commentLists,
	})
}
