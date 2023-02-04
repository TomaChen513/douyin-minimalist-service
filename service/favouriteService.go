package service

import (
	"github.com/RaymondCode/simple-demo/model"
)

type Favor struct {
	Id            int64  `json:"id,omitempty"`
	Author        User   `json:"author,omitempty"`
	PlayUrl       string `json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
	Title         string `json:"title,omitempty"`
}

type FavorService interface {
	// 点赞操作
	FavoriteAction(userId, videoId int64, actionType string) bool
	GetFavouriteList(userId int64) ([]Favor, error)
	// //IsFavorite 根据当前视频id判断是否点赞了该视频。
	// IsFavourite(videoId int64, userId int64) (bool, error)
	//FavouriteCount 根据当前视频id获取当前视频点赞数量。
	FavouriteCount(videoId int64) (int64, error)
	// // GetUserByIdWithCurId 已登录(curID)情况下,根据user_id获得User对象
	// GetUserByIdWithCurId(id int64, curId int64) (User, error)
}

type FavorServiceImpl struct {
	UserService
	VideoService
}

func (fvsi *FavorServiceImpl) FavoriteAction(userId, videoId int64, actionType string) bool {
	if actionType == "1" {
		// 执行点赞
		return model.InsertFavourite(model.Like{UserId: userId, VideoId: videoId})
	} else if actionType == "2" {
		// 执行取消点赞
		return model.DeleteFavourite(model.Like{UserId: userId, VideoId: videoId})
	}
	return false
}

func (fvsi *FavorServiceImpl) GetFavouriteList(userId int64) ([]Favor, error) {
	// 查找所有该用户喜欢的视频下标
	favouriteVideosId, _ := model.SelectVideosByUserId(userId)
	// 根据视频下标获取视频信息
	videos, _ := model.GetVideosById(favouriteVideosId)

	favors := make([]Favor, len(videos))

	for i := 0; i < len(videos); i++ {
		favors[i].Id = videos[i].Id
		author, _ := fvsi.GetUserById(videos[i].AuthorId)
		favors[i].Author = author
		favors[i].PlayUrl = videos[i].PlayUrl
		favors[i].CoverUrl = videos[i].CoverUrl
		// 还未实现，使用默认值
		favors[i].FavoriteCount = 1
		favors[i].CommentCount = 1
		favors[i].IsFavorite = true
		favors[i].Title = "testTitle"
	}
	return favors, nil
}

func (fvsi *FavorServiceImpl) FavoriteCount(videoId int64) (int64,error){
	videos,err:=model.SelectLikesByVideoId(videoId)
	if err!=nil {
		return -1,err
	}
	return len(videos),nil
}
