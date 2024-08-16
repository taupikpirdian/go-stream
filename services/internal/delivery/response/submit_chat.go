package response

import (
	"time"

	chatPb "cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/pb"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/entity"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewSubmitChatResponse(u *entity.ChatBotSubmit) *chatPb.SubmitQuestionResponse {
	submitDate := u.SubmitDateTime.Add(time.Hour * 7)
	uPb := chatPb.SubmitQuestionData{
		ConversationId: u.ConversationId,
		MessageId:      u.MessageId,
		TaskId:         u.TaskId,
		SubmitDate:     timestamppb.New(submitDate),
		Answer:         u.Answer,
		Feedback:       "NOT_DEFINE",
		RequestId:      u.RequestId,
		Status:         u.Status,
	}
	return &chatPb.SubmitQuestionResponse{
		Data: &uPb,
	}
}
