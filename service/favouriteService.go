package service

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/RaymondCode/simple-demo/lib"
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
	// 根据用户id获得点赞列表
	GetFavouriteList(userId int64) ([]Video, error)
	// 根据视频id获取视频点赞数量。
	FavouriteCount(videoId int64) (int64, error)

	// //IsFavorite 根据当前视频id判断是否点赞了该视频。
	// IsFavourite(videoId int64, userId int64) (bool, error)
	// // GetUserByIdWithCurId 已登录(curID)情况下,根据user_id获得User对象
	// GetUserByIdWithCurId(id int64, curId int64) (User, error)
}

type FavorServiceImpl struct {
	UserService
	VideoService
}

// 重复点赞问题
// 点赞操作
func (fvsi *FavorServiceImpl) FavoriteAction(userId, videoId int64, actionType string) bool {
	if actionType == "1" {
		// 执行点赞
		return fvsi.likeAction(userId, videoId)
	} else if actionType == "2" {
		// 执行取消点赞
		return fvsi.unLikeAction(userId, videoId)
	}
	return false
}

// ===============待编写测试文件====================
// 根据用户id获得点赞列表
// 1. 查询缓存  2. 缓存存在，获得所有对应的videoId，若不存在，进入数据库查询
// 3. 获得videoId后，并发的去数据库获得video的信息
func (fvsi *FavorServiceImpl) GetFavouriteList(userId int64) ([]Video, error) {
	uId := strconv.FormatInt(userId, 10)
	// 开辟video切片空间
	favoriteVideoList := make([]Video, 0)
	// 若redis中存在userId的缓存
	if n, err := lib.RdbLikeUserId.Exists(lib.Ctx, uId).Result(); n > 0 {
		if err != nil {
			log.Printf("方法:GetFavouriteList：查询redis缓存(userId——>videoId)时出错:%v", err)
			return nil, err
		}
		// 获取全部videoId
		videoIds, err1 := lib.RdbLikeUserId.SMembers(lib.Ctx, uId).Result()
		if err1 != nil {
			log.Printf("方法:GetFavouriteList：执行redis SMEMBER出错:%v", err)
		}

		// 因为第一次加载时有一个-1，所以需要判断
		videoLength := len(videoIds) - 1
		if videoLength == 0 {
			return favoriteVideoList, nil
		}

		// 使用协程逐个获取video对象
		var wg sync.WaitGroup
		wg.Add(videoLength)

		for i := 0; i <= videoLength; i++ {
			videoId, _ := strconv.ParseInt(videoIds[i], 10, 64)
			// 遇到默认值跳过
			if videoId == -1 {
				continue
			}
			// 进入协程逐个获取video
			go fvsi.addVideoList(videoId, userId, &favoriteVideoList, &wg)
		}
		// 等待所有协程执行完毕
		wg.Wait()
		return favoriteVideoList, nil
	}
	// 若不存在则直接数据库并发查找并存入缓存

	//声明userId空间
	if _, err := lib.RdbLikeUserId.SAdd(lib.Ctx, uId, -1).Result(); err != nil {
		log.Printf("方法:GetFavouriteList RedisLikeUserId add value失败")
		lib.RdbLikeUserId.Del(lib.Ctx, uId)
		return nil, err
	}

	// 设置过期时间
	_, err := lib.RdbLikeUserId.Expire(lib.Ctx, uId, time.Duration(24*time.Hour)).Result()
	if err != nil {
		log.Printf("方法:GetFavouriteList RedisLikeUserId 设置有效期失败")
		lib.RdbLikeUserId.Del(lib.Ctx, uId)
		return nil, err
	}

	// 查找所有该用户喜欢的视频下标
	favouriteVideosId, err := model.SelectVideosByUserId(userId)
	if err != nil {
		log.Printf("方法:GetFavouriteList model.SelectVideosByUserId失败：%v", err)
		lib.RdbLikeUserId.Del(lib.Ctx, uId)
		return nil, err
	}

	//遍历videoIdList,添加进key的集合中，若失败，删除key
	//防止脏读，保证redis与mysql数据一致性
	for _, vId := range favouriteVideosId {
		if _, err := lib.RdbLikeUserId.SAdd(lib.Ctx, uId, vId).Result(); err != nil {
			log.Printf("方法:GetFavouriteList RedisLikeUserId add value失败")
			lib.RdbLikeUserId.Del(lib.Ctx, uId)
			return nil, err
		}
	}

	videoLength := len(favouriteVideosId)
	// 使用协程逐个获取video对象,这里采用的video id是直接从数据库获取的
	var wg sync.WaitGroup
	wg.Add(videoLength)

	for i := 0; i < videoLength; i++ {
		// 进入协程逐个获取video
		go fvsi.addVideoList(favouriteVideosId[i], userId, &favoriteVideoList, &wg)
	}
	// 等待所有协程执行完毕
	wg.Wait()
	return favoriteVideoList, nil
}

// ===============待编写测试文件====================

// 点赞操作
// 1. 对缓存进行添加操作，若无缓存  同时开协程修改视频点赞数量
func (fvsi *FavorServiceImpl) likeAction(userId, videoId int64) bool {

	return false
}
func (fvsi *FavorServiceImpl) unLikeAction(userId, videoId int64) bool {
	return false
}

// ===============待编写测试文件====================
// 根据videoId和userId从数据库获取一条video信息
func (fvsi *FavorServiceImpl) addVideoList(videoId, userId int64, videoList *[]Video, wg *sync.WaitGroup) {
	defer wg.Done()
	video, err := fvsi.GetVideo(videoId, userId)
	if err != nil {
		log.Printf("执行VideoService:GetVideo出错:%v", err)
	}
	*videoList = append(*videoList, video)
}

// ===============待编写测试文件====================

// ===============待编写测试文件====================
// 根据video id获得视频点赞数量,-1表示查询的数据有误
// 1. 判断缓存，若无缓存，从数据库读取并存入缓存
func (fvsi *FavorServiceImpl) FavouriteCount(videoId int64) (int64, error) {
	strVideoId := strconv.FormatInt(videoId, 10)
	// 若点赞数存在缓存中
	if n, err := lib.RdbLikeVideoCount.Exists(lib.Ctx, strVideoId).Result(); n > 0 {
		// 出现查询存在key，但是失败的情况
		if err != nil {
			log.Printf("方法:FavouriteCount RedisLikeVideoCount query key失败：%v", err)
			return -1, err
		}
		strCount, err1 := lib.RdbLikeVideoCount.Get(lib.Ctx, strVideoId).Result()
		if err1 != nil {
			log.Printf("方法:FavouriteCount RedisLikeVideoCount key存储value有误：%v", err1)
			return -1, err
		}
		count, _ := strconv.ParseInt(strCount, 10, 64)
		return count, nil
	}

	// 若不在缓存中则直接读取数据库并设置进缓存，同时设置失效时间
	// 数据库获取点赞量
	count, err := model.CountLikesByVideoId(videoId)
	if err != nil {
		log.Printf("方法:FavouriteCount model.CountLikesByVideoId：%v", err)
		return -1, err
	}

	// 存入redis并设置一天的过期时间，这里怕有并发安全问题
	_,err1:=lib.RdbLikeVideoCount.Set(lib.Ctx,strVideoId,count,24*time.Hour).Result()
	if err1!=nil {
		log.Printf("方法:FavouriteCount 存入redis出现错误:%v", err1)
		return -1,err1
	}

	return count, nil
}
// ===============待编写测试文件====================
