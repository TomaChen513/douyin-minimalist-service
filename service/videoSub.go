package service

type VideoSub struct {
}

func (vs VideoSub) GetUserByIdWithCurId(id int64, curId int64) (User, error) {
	var user User
	return user, nil
}

func (vs VideoSub) FavouriteCount(videoId int64) (int64, error) {
	return 3, nil
}
