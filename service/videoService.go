package service

import "github.com/RaymondCode/simple-demo/model"


type VideoService interface {
	GetVideosByUser(userId int64) ([]Video,error)
	InsertVideo(video Video) bool
}

type Video struct {
	Id            int64  `gorm:"primarykey" json:"id,omitempty"`
	Author        User   `json:"author"`
	PlayUrl       string `gorm:"type:varchar(255)" json:"play_url,omitempty"`
	CoverUrl      string `gorm:"type:varchar(255)" json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
	Title string `gorm:"type:varchar(55)" json:"title,omitempty"`
}

type VideoServiceImpl struct {
	UserServiceImpl
}

// 由于在获得用户投稿信息中只有用户id，还需要根据id获取用户信息
// 为避免controller层逻辑过于复杂，因此在service层中实现
// 这个地方性能可以优化
func (vsi *VideoServiceImpl)GetVideosByUser(userId int64) ([]Video,error){
	tableVideos,err:=model.GetVideosByUserId(userId)
	videos:=make([]Video,len(tableVideos))

	if err!=nil {
		return videos,err
	}

	for i := 0; i < len(tableVideos); i++ {
		user,_:=vsi.GetUserById(tableVideos[i].Id)
		videos[i].Author=user
		videos[i].Id=tableVideos[i].Id
		videos[i].PlayUrl=tableVideos[i].PlayUrl
		videos[i].CoverUrl=tableVideos[i].CoverUrl
		videos[i].FavoriteCount=tableVideos[i].FavoriteCount
		videos[i].CommentCount=tableVideos[i].CommentCount
		videos[i].IsFavorite=tableVideos[i].IsFavorite
		videos[i].Title=tableVideos[i].Title
	}
	return videos,nil
}

func (vsi *VideoServiceImpl)InsertVideo(video Video)bool{
	tableVideo:=model.TableVideo{
		Id: video.Id,
		Author: video.Author.Id,
		PlayUrl: video.PlayUrl,
		CoverUrl: video.CoverUrl,
		FavoriteCount: video.FavoriteCount,
		CommentCount: video.CommentCount,
		IsFavorite: video.IsFavorite,
		Title: video.Title,
	}
	return model.InsertVideo(tableVideo)
}