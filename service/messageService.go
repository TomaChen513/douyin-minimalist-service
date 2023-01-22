package service

import "github.com/RaymondCode/simple-demo/model"

type MessageService interface {
	SendMessage(uId, tUId int64, content string) bool
	SelectMessageByUserId(uId, tUId int64) ([]Message, error)
}

type Message struct {
	Id         int64  `json:"id,omitempty"`
	Content    string `json:"content,omitempty"`
	CreateTime string `json:"create_time,omitempty"`
}

type MessageServiceImpl struct {
}

func (msi *MessageServiceImpl) SendMessage(uId, tUId int64, content string) bool {
	return model.InsertMessage(model.TableMessage{UserId: uId, ToUserId: tUId,Content: content})
}


// 没有考虑相互，只是写一个demo，后续改用其他技术实现
func (msi *MessageServiceImpl) SelectMessageByUserId(uId, tUId int64) ([]Message, error) {
	tableMessages, err := model.SelectMessagesByUserId(uId, tUId)
	messages := make([]Message, len(tableMessages))
	if err != nil {
		return messages, err
	}
	for i := 0; i < len(tableMessages); i++ {
		messages[i] = Message{Id: tableMessages[i].Id,
			Content:    tableMessages[i].Content,
			CreateTime: tableMessages[i].CreateTime.String()}
	}
	return messages, nil
}
