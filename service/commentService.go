package service

import (
	"log"
	"strconv"
	"time"

	"github.com/RaymondCode/simple-demo/lib"
	"github.com/RaymondCode/simple-demo/model"
)

type CommentService interface {
	// 根据评论id删除评论，不缺定前端是否会判断该评论是否属于该用户，这里假设传来的评论id属于对应用户
	DeleteComment(commentId, userId, videoId int64) bool
	// 发表评论
	InsertComment(userId, videoId int64, content string) (Comment, error)
	// 按时间倒序获得评论
	GetCommentListDecByTime(videoId int64) ([]Comment, error)
	// 根据videoId获取视频评论数量
	CountByVideoId(videoId int64) (int64, error)
}

type Comment struct {
	Id         int64  `json:"id,omitempty"`
	User       User   `json:"user"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}

type CommentServiceImpl struct {
	UserService
}

// 根据主键删除评论
// bug:默认一条视频用户只能发一个评论，存在错误需要修复
// 由于评论一定是一对一的，不可能出现同时两个人删除同一条评论的情况。
// 1. 查询缓存，是否存在，存在删除缓存，推入消息队列。若不存在，直接推入消息队列,同时count缓存减1
func (csi *CommentServiceImpl) DeleteComment(commentId, userId, videoId int64) bool {
	strUserId := strconv.FormatInt(userId, 10)
	strVideoId := strconv.FormatInt(videoId, 10)
	commentKey := strVideoId + "_" + strUserId
	// RdbCommentId
	// 若存在缓存
	if n, err := lib.RdbCommentid.Exists(lib.Ctx, commentKey).Result(); n > 0 {
		if err != nil {
			log.Printf("方法:DeleteComment RedisCommentId query key失败:%v", err)
			return false
		}

		// 删除缓存
		if _, err1 := lib.RdbCommentid.Del(lib.Ctx, commentKey).Result(); err1 != nil {
			log.Printf("方法:DeleteComment RedisCommentId 删除 key失败:%v", err)
		}
		lib.RdbCommentVideoId.SRem(lib.Ctx,strVideoId,strUserId)
	}

	// RdbCommentCount
	// 若不存在，则先新建缓存
	if _, err := lib.RdbCommentCount.Exists(lib.Ctx, strVideoId).Result(); err != nil {
		log.Printf("方法:DeleteComment RedisCommentCount 不存在：%v", err)
		_, err1 := csi.CountByVideoId(videoId)
		// 若获取评论数出现错误则回退之前操作
		if err1 != nil {
			log.Printf("方法:DeleteComment RedisCommentCount 获取点赞错误：%v", err1)
			return false
		}
	}
	// 缓存值减1
	_, err := lib.RdbCommentCount.Decr(lib.Ctx, strVideoId).Result()
	if err != nil {
		log.Printf("方法:DeleteComment RedisCommentCount 自减错误：%v", err)
		return false
	}
	// 推入消息队列
	// 消息队列
	model.DeleteComment(commentId)
	return true
}

// 插入新评论
// 1. 对commentId缓存进行添加操作
func (csi *CommentServiceImpl) InsertComment(userId, videoId int64, content string) (Comment, error) {
	strUserId := strconv.FormatInt(userId, 10)
	strVideoId := strconv.FormatInt(videoId, 10)
	commentKey := strVideoId + "_" + strUserId

	// RdbCommentid 加入缓存
	if _, err1 := lib.RdbCommentid.Set(lib.Ctx, commentKey, content, 24*time.Hour).Result(); err1 != nil {
		log.Printf("方法:InsertComment RdbCommentid add value失败:%v", err1)
		return Comment{}, err1
	}

	// RdbCommentVideoId 判断缓存
	// 若存在缓存
	if n, err := lib.RdbCommentVideoId.Exists(lib.Ctx, strVideoId).Result(); n > 0 {
		if err != nil {
			return Comment{},err
		}

		// 将userId加入缓存
		if _, err1 := lib.RdbCommentVideoId.SAdd(lib.Ctx, strVideoId, strUserId).Result(); err1 != nil {
			return Comment{},err1
		}
	} else {
		if _, err := lib.RdbCommentVideoId.SAdd(lib.Ctx, strVideoId, -1).Result(); err != nil {
			lib.RdbCommentid.Del(lib.Ctx, commentKey)
			return Comment{},err
		}
		// 设置过期时间
		_, err := lib.RdbCommentVideoId.Expire(lib.Ctx, strVideoId, 24*time.Hour).Result()
		if err != nil {
			lib.RdbCommentid.Del(lib.Ctx, commentKey)
			return Comment{},err
		}
		// 查询数据库获得videoId所有评论userid，存入redis
		userIds, err1 := model.GetUserIdByVideoId(videoId)
		if err1 != nil {
			lib.RdbCommentid.Del(lib.Ctx, commentKey)
		}

		// 将所有userId存入缓存，若失败删除key，并返回，为了防止脏读
		for _, userId := range userIds {
			if _, err := lib.RdbCommentVideoId.SAdd(lib.Ctx, strUserId, userId).Result(); err != nil {
				lib.RdbCommentid.Del(lib.Ctx, commentKey)
				return Comment{},err
			}
		}
		if _, err2 := lib.RdbCommentVideoId.SAdd(lib.Ctx, strVideoId, strUserId).Result(); err2 != nil {
			return Comment{},err2
		}
	}

	// RdbLikeVideoCount 先进行自增，若自增错误则删除上文插入的key
	// 若不存在，则先新建缓存
	if _, err := lib.RdbCommentCount.Exists(lib.Ctx, strVideoId).Result(); err != nil {
		log.Printf("方法:InsertCommentn RedisCommentCount 不存在：%v", err)
		_, err1 := csi.CountByVideoId(videoId)
		// 若获取评论数出现错误则回退之前操作
		if err1 != nil {
			log.Printf("方法:InsertComment RedisCommentCount 获取评论错误：%v", err1)
			_, err2 := lib.RdbCommentid.Del(lib.Ctx, commentKey).Result()
			if err2 != nil {
				log.Printf("方法:InsertComment 移除元素错误：%v", err1)
				return Comment{}, nil
			}
			lib.RdbCommentVideoId.SRem(lib.Ctx,strUserId).Result()
			return Comment{}, nil
		}
	}
	// 缓存值加一
	_, err := lib.RdbCommentCount.Incr(lib.Ctx, commentKey).Result()
	if err != nil {
		log.Printf("方法:InsertComment 自增错误：%v", err)
		return Comment{}, nil
	}
	// 若无错误，加入消息队列操作数据库
	// 若数据库操作失败，缓存已经存在，不影响
	// rabbitmq
	user, _ := csi.GetUserById(userId)
	createTime := time.Now()
	tableComment := model.Comment{UserId: userId, VideoId: videoId, CommentText: content, CreateDate: createTime}
	if !model.InsertComment(tableComment) {
		return Comment{}, nil
	}
	return Comment{User: user, Content: content, CreateDate: createTime.UTC().String()}, nil
	//
	//
}

// 按时间倒序获得视频评论
// 1.查找缓存获得userId 2.查找缓存获得comment
func (csi *CommentServiceImpl) GetCommentListDecByTime(videoId int64) ([]Comment, error) {
	strVideoId := strconv.FormatInt(videoId, 10)

	// RdbCommentVideoId
	// 若有缓存
	if n,err:=lib.RdbCommentVideoId.Exists(lib.Ctx,strVideoId).Result() ;n>0{
		if err!=nil {
			return []Comment{},err
		}
		// 未完善
		lib.RdbCommentVideoId.SMembers(lib.Ctx,strVideoId).Result()
		return []Comment{},nil

	}else{
		// 数据库查找
		tableComments, err := model.GetCommentListDecByTime(videoId)
		commentList := make([]Comment, len(tableComments))
		if err != nil {
			return commentList, nil
		}
		for i := 0; i < len(tableComments); i++ {
			user, _ := csi.GetUserById(tableComments[i].UserId)
			commentList[i] = Comment{
				Id:         tableComments[i].Id,
				User:       user,
				Content:    tableComments[i].CommentText,
				CreateDate: tableComments[i].CreateDate.String(),
			}
		}
		return commentList, nil
	}
}

// 根据videoId获取视频评论数量
func (csi *CommentServiceImpl) CountByVideoId(videoId int64) (int64, error) {
	strVideoId := strconv.FormatInt(videoId, 10)
	// 若评论数存在缓存中
	if n, err := lib.RdbCommentCount.Exists(lib.Ctx, strVideoId).Result(); n > 0 {
		// 出现查询存在key，但是失败的情况
		if err != nil {
			log.Printf("方法:CountByVideoId RedisCommentCount query key失败:%v", err)
			return -1, err
		}
		strCount, err1 := lib.RdbCommentCount.Get(lib.Ctx, strVideoId).Result()
		if err1 != nil {
			log.Printf("方法:CountByVideoId RedisCommentCount key存储value有误：%v", err1)
			return -1, err
		}
		count, _ := strconv.ParseInt(strCount, 10, 64)
		return count, nil
	}

	// 若不在缓存中则直接读取数据库并设置进缓存，同时设置失效时间
	// 数据库获取评论量
	count, err := model.CountCommentsByVideoId(videoId)
	if err != nil {
		log.Printf("方法:CountByVideoId CountCommentsByVideoId%v", err)
		return -1, err
	}

	// 存入redis并设置一天的过期时间，这里怕有并发安全问题
	_, err1 := lib.RdbCommentCount.Set(lib.Ctx, strVideoId, count, 24*time.Hour).Result()
	if err1 != nil {
		log.Printf("方法:CountByVideoId 存入redis出现错误:%v", err1)
		return -1, err1
	}

	return count, nil
}
