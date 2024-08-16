package service

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/pkg/utils"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/entity"
	utilPkg "cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/utils"
)

func ManageResponseAndError(response *http.Response, err error) ([]byte, error) {
	// check request success/failure
	if err != nil {
		return nil, errors.New(utilPkg.ClassifyNetworkError(err))
	}

	// read response body
	res, err := io.ReadAll(response.Body)
	// leverage defer stack to defer closing of response body read operation
	// this will defer until this function is ready to return

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	if err != nil {
		return res, errors.New("invalid response body")
	}
	// check response status
	if response.StatusCode == 403 {
		return res, errors.New("not authorized")
	}

	return res, nil
}

func (u *ResponseFeedbackAnswer) CriteriaSuccessFeedbackAnswer() error {
	if u.Result != "success" {
		return errors.New("criteria not success")
	}
	return nil
}

func (u *ResponseStopAnswer) CriteriaSuccessStopAnswer() error {
	if u.Result != "success" {
		return errors.New("criteria not success")
	}
	return nil
}

func MapDetailAnswer(res *ResponseDetailAnswer, conversationId string) *entity.DetailAnswer {
	chatList := MapDetailChatList(res.Data)
	return &entity.DetailAnswer{
		ConversationId: conversationId,
		ChatDate:       utils.UnixToTime(res.MetaData.Conversation.CreatedAt),
		ChatList:       chatList,
	}
}

func MapDetailChatList(res []DataResponseDetailAnswer) []*entity.DetailChatList {
	chatList := make([]*entity.DetailChatList, 0)
	for _, v := range res {
		chatList = append(chatList, &entity.DetailChatList{
			Id:       v.Id,
			Question: v.Query,
			Answers: []*entity.DetailChatListAnswer{
				{
					GenerateTime: utils.UnixToTime(v.CreatedAt),
					Answer:       v.Answer,
					Feedback:     strings.ToUpper(v.FeedBack.Rating),
				},
			},
		})
	}
	return chatList
}

func MapConversation(res *ResponseConversation) *entity.ListConversation {
	conversationList := MapListConversation(res.Data)
	return &entity.ListConversation{
		Data: conversationList,
		Meta: entity.MetaData{
			HasMore: res.HasMore,
		},
	}
}

func MapListConversation(res []DataResponseConversation) []*entity.DataConvesation {
	converationList := make([]*entity.DataConvesation, 0)
	for _, v := range res {
		converationList = append(converationList, &entity.DataConvesation{
			Id:       v.Id,
			Title:    v.Name,
			ChatDate: utils.EpochToTime(int(v.CreatedAt)),
		})
	}
	return converationList
}
