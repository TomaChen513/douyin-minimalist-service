package service

import (
	"log"
	"time"

	"github.com/RaymondCode/simple-demo/lib"
	"github.com/RaymondCode/simple-demo/model"
)

type VideoServiceImpl struct {
	UserService
	FavorService
	CommentService
}

// Feed
// 通过传入时间戳，当前用户的id，返回对应的视频数组，以及视频数组中最早的发布时间
// 获取视频数组大小是可以控制的，在config中的videoCount变量
func (videoService VideoServiceImpl) Feed(lastTime time.Time, userId int64) ([]Video, time.Time, error) {
	//创建对应返回视频的切片数组，提前将切片的容量设置好，可以减少切片扩容的性能
	videos := make([]Video, 0, lib.VideoCount)
	//根据传入的时间，获得传入时间前n个视频，可以通过config.videoCount来控制
	tableVideos, err := model.GetVideosByLastTime(lastTime)
	if err != nil {
		log.Printf("方法GetVideosByLastTime(lastTime) 失败：%v", err)
		return nil, time.Time{}, err
	}
	log.Printf("方法GetVideosByLastTime(lastTime) 成功：%v", tableVideos)
	//将数据通过copyVideos进行处理，在拷贝的过程中对数据进行组装
	err = videoService.copyVideos(&videos, &tableVideos, userId)
	if err != nil {
		log.Printf("方法videoService.copyVideos(&videos, &tableVideos, userId) 失败：%v", err)
		return nil, time.Time{}, err
	}
	log.Printf("方法videoService.copyVideos(&videos, &tableVideos, userId) 成功")
	//返回数据，同时获得视频中最早的时间返回
	return videos, tableVideos[len(tableVideos)-1].PublishTime, nil
}

// List
// 通过userId来查询对应用户发布的视频，并返回对应的视频数组
func (videoService *VideoServiceImpl) List(userId int64, curId int64) ([]Video, error) {
	//依据用户id查询所有的视频，获取视频列表
	data, err := model.GetVideosByAuthorId(userId)
	if err != nil {
		log.Printf("方法GetVideosByAuthorId(%v)失败:%v", userId, err)
		return nil, err
	}
	log.Printf("方法GetVideosByAuthorId(%v)成功:%v", userId, data)
	//提前定义好切片长度
	result := make([]Video, 0, len(data))
	//调用拷贝方法，将数据进行转换
	err = videoService.copyVideos(&result, &data, curId)
	if err != nil {
		log.Printf("方法videoService.copyVideos(&result, &data, %v)失败:%v", userId, err)
		return nil, err
	}
	//如果数据没有问题，则直接返回
	return result, nil
}

// 该方法可以将数据进行拷贝和转换，并从其他方法获取对应的数据
func (videoService *VideoServiceImpl) copyVideos(result *[]Video, data *[]model.Video, userId int64) error {
	for _, temp := range *data {
		var video Video
		//将video进行组装，添加想要的信息,插入从数据库中查到的数据
		videoService.creatVideo(&video, &temp, userId)
		*result = append(*result, video)
	}
	return nil
}

// 将video进行组装，添加想要的信息,插入从数据库中查到的数据
func (videoService *VideoServiceImpl) creatVideo(video *Video, data *model.Video, userId int64) {
	//建立协程组，当这一组的携程全部完成后，才会结束本方法
	// var wg sync.WaitGroup
	// wg.Add(4)
	var err error
	video.Video = *data
	//插入Author，这里需要将视频的发布者和当前登录的用户传入，才能正确获得isFollow，
	//如果出现错误，不能直接返回失败，将默认值返回，保证稳定
	go func() {
		video.Author, err = videoService.GetUserByIdWithCurId(data.AuthorId, userId)
		if err != nil {
			log.Printf("方法videoService.GetUserByIdWithCurId(data.AuthorId, userId) 失败：%v", err)
		} else {
			log.Printf("方法videoService.GetUserByIdWithCurId(data.AuthorId, userId) 成功")
		}
		// wg.Done()
	}()

	// //插入点赞数量，同上所示，不将nil直接向上返回，数据没有就算了，给一个默认就行了
	// go func() {
	// 	video.FavoriteCount, err = videoService.FavouriteCount(data.Id)
	// 	if err != nil {
	// 		log.Printf("方法videoService.FavouriteCount(data.ID) 失败：%v", err)
	// 	} else {
	// 		log.Printf("方法videoService.FavouriteCount(data.ID) 成功")
	// 	}
	// 	wg.Done()
	// }()

	// //获取该视屏的评论数字
	// go func() {
	// 	video.CommentCount, err = videoService.CountFromVideoId(data.Id)
	// 	if err != nil {
	// 		log.Printf("方法videoService.CountFromVideoId(data.ID) 失败：%v", err)
	// 	} else {
	// 		log.Printf("方法videoService.CountFromVideoId(data.ID) 成功")
	// 	}
	// 	wg.Done()
	// }()

	// //获取当前用户是否点赞了该视频
	// go func() {
	// 	video.IsFavorite, err = videoService.IsFavourite(video.Id, userId)
	// 	if err != nil {
	// 		log.Printf("方法videoService.IsFavourit(video.Id, userId) 失败：%v", err)
	// 	} else {
	// 		log.Printf("方法videoService.IsFavourit(video.Id, userId) 成功")
	// 	}
	// 	wg.Done()
	// }()

	// wg.Wait()
}

// GetVideo
// 传入视频id获得对应的视频对象，注意还需要传入当前登录用户id
func (videoService *VideoServiceImpl) GetVideo(videoId int64, userId int64) (Video, error) {
	//初始化video对象
	var video Video

	//从数据库中查询数据，如果查询不到数据，就直接失败返回，后续流程就不需要执行了
	data, err := model.GetVideoByVideoId(videoId)
	if err != nil {
		log.Printf("方法dao.GetVideoByVideoId(videoId) 失败：%v", err)
		return video, err
	} else {
		log.Printf("方法dao.GetVideoByVideoId(videoId) 成功")
	}

	//插入从数据库中查到的数据
	videoService.creatVideo(&video, &data, userId)
	return video, nil
}

// GetVideoIdList
// 通过一个作者id，返回该用户发布的视频id切片数组
func (videoService *VideoServiceImpl) GetVideoIdList(authorId int64) ([]int64, error) {
	//直接调用dao层方法获取id即可
	id, err := model.GetVideoIdsByAuthorId(authorId)
	if err != nil {
		log.Printf("方法dao.GetVideoIdsByAuthorId(%v) 失败：%v", authorId, err)
		return nil, err
	} else {
		log.Printf("方法dao.GetVideoIdsByAuthorId(%v) 成功", authorId)
	}
	return id, nil
}
