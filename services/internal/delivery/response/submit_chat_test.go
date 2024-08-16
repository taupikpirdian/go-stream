package response_test

import (
	"testing"

	chatPb "cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/pb"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/entity"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/internal/delivery/response"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestNewSubmitChatResponse(t *testing.T) {
	type args struct {
		u *entity.ChatBotSubmit
	}
	tests := []struct {
		name string
		args args
		want *chatPb.SubmitQuestionResponse
	}{
		{
			name: "success",
			args: args{
				u: &entity.ChatBotSubmit{
					ConversationId: "39312c44-c579-427a-b513-33ae492729dc",
					MessageId:      "39312c44-c579-427a-b513-33ae492729d1",
					TaskId:         "39312c44-c579-427a-b513-33ae492729d2",
					SubmitDate:     1713952679,
					Answer:         "In publishing and graphic design, Lorem ipsum is a placeholder text commonly used to demonstrate the visual form of a document or a typeface without relying on meaningful content.",
					Feedback:       "",
					RequestId:      "39312c44-c579-427a-b513-33ae492729d3",
					Status:         "pending",
				},
			},
			want: &chatPb.SubmitQuestionResponse{
				Data: &chatPb.SubmitQuestionData{
					ConversationId: "39312c44-c579-427a-b513-33ae492729dc",
					MessageId:      "39312c44-c579-427a-b513-33ae492729d1",
					TaskId:         "39312c44-c579-427a-b513-33ae492729d2",
					SubmitDate:     &timestamppb.Timestamp{},
					Answer:         "In publishing and graphic design, Lorem ipsum is a placeholder text commonly used to demonstrate the visual form of a document or a typeface without relying on meaningful content.",
					Feedback:       "",
					RequestId:      "39312c44-c579-427a-b513-33ae492729d3",
					Status:         "pending",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := response.NewSubmitChatResponse(tt.args.u)
			assert.Equal(t, tt.want.Data.ConversationId, got.Data.ConversationId)
			assert.Equal(t, tt.want.Data.Answer, got.Data.Answer)
		})
	}
}
