package service

import (
	"time"

	"github.com/RaymondCode/simple-demo/model"
)

type CommentService interface {
	DeleteComment(commentId int64) bool
	InsertComment(userId,videoId int64,content string)(Comment,error)
	GetCommentListDecByTime(videoId int64) ([]Comment,error)
}

type Comment struct {
	Id         int64  `json:"id,omitempty"`
	User       User   `json:"user"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}

type CommentServiceImpl struct {
	UserServiceImpl
}

// 根据主键删除评论
func (csi *CommentServiceImpl)DeleteComment(commentId int64) bool{
	return model.DeleteComment(commentId)
}

// 插入新评论
func (csi *CommentServiceImpl)InsertComment(userId,videoId int64,content string)(Comment,error){
	user,_:=csi.GetUserById(userId)
	createTime:=time.Now()
	tableComment:=model.TableComment{UserId: userId,VideoId: videoId,Content: content,CreateDate: createTime}
	if !model.InsertComment(tableComment){
		return Comment{},nil
	}
	return Comment{User: user,Content: content,CreateDate: createTime.UTC().String()},nil
}

// 按时间倒序获得视频评论
func (csi *CommentServiceImpl)GetCommentListDecByTime(videoId int64)([]Comment,error){
	tableComments,err:=model.GetCommentListDecByTime(videoId)
	commentList:=make([]Comment,len(tableComments))
	if err!=nil {
		return commentList,nil
	}
	for i := 0; i < len(tableComments); i++ {
		user,_:=csi.GetUserById(tableComments[i].UserId)
		commentList[i]=Comment{
			Id: tableComments[i].Id,
			User: user,
			Content: tableComments[i].Content,
			CreateDate: tableComments[i].CreateDate.String(),
		}
	}
	return commentList,nil
}