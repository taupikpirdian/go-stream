package response

import (
	chatPb "cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/pb"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/entity"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewListConversationResponse(u *entity.ListConversation) *chatPb.GetConversationListResponse {
	listConversation := make([]*chatPb.DataGetConversationList, 0)
	for _, v := range u.Data {
		listConversation = append(listConversation, &chatPb.DataGetConversationList{
			Id:       v.Id,
			ChatDate: timestamppb.New(v.ChatDate),
			Title:    v.Title,
		})
	}
	return &chatPb.GetConversationListResponse{
		Data: listConversation,
		Meta: &chatPb.MetaData{
			HasMore: u.Meta.HasMore,
		},
	}
}
