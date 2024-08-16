package response

import (
	chatPb "cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/pb"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/entity"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewDetailAnswerResponse(u *entity.DetailAnswer) *chatPb.GetDetailResponse {
	chatList := SetChatList(u.ChatList)
	uPb := chatPb.DataGetDetailResponse{
		ConversationId: u.ConversationId,
		ChatDate:       timestamppb.New(u.ChatDate),
		ChatList:       chatList,
	}
	return &chatPb.GetDetailResponse{
		Data: &uPb,
	}
}

func SetChatList(d []*entity.DetailChatList) []*chatPb.DataChatList {
	chatList := make([]*chatPb.DataChatList, 0)
	for _, v := range d {
		answers := SetAnswerList(v.Answers)
		chatList = append(chatList, &chatPb.DataChatList{
			Id:       v.Id,
			Question: v.Question,
			Answers:  answers,
		})
	}
	return chatList
}

func SetAnswerList(d []*entity.DetailChatListAnswer) []*chatPb.DataAnswer {
	answerList := make([]*chatPb.DataAnswer, 0)
	for _, v := range d {
		answerList = append(answerList, &chatPb.DataAnswer{
			GenerateTime: timestamppb.New(v.GenerateTime),
			Answer:       v.Answer,
			Feedback:     v.Feedback,
		})
	}
	return answerList
}
