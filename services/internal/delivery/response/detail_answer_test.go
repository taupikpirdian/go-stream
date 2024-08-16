package response_test

import (
	"testing"
	"time"

	chatPb "cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/pb"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/entity"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/internal/delivery/response"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestNewDetailAnswerResponse(t *testing.T) {
	timeNow := time.Now()
	type args struct {
		u *entity.DetailAnswer
	}
	tests := []struct {
		name string
		args args
		want *chatPb.GetDetailResponse
	}{
		{
			name: "success - mapping data to response",
			args: args{
				u: &entity.DetailAnswer{
					ConversationId: "39312c44-c579-427a-b513-33ae492729dc",
					ChatDate:       timeNow,
					ChatList: []*entity.DetailChatList{
						{
							Id:       "39312c44-c579-427a-b513-33ae492729d1",
							Question: "Tes Question",
							Answers: []*entity.DetailChatListAnswer{
								{
									GenerateTime: timeNow,
									Answer:       "Tes Answer",
									Feedback:     "LIKE",
								},
							},
						},
					},
				},
			},
			want: &chatPb.GetDetailResponse{
				Data: &chatPb.DataGetDetailResponse{
					ConversationId: "39312c44-c579-427a-b513-33ae492729dc",
					ChatDate:       timestamppb.New(timeNow),
					ChatList: []*chatPb.DataChatList{
						{
							Id:       "39312c44-c579-427a-b513-33ae492729d1",
							Question: "Tes Question",
							Answers: []*chatPb.DataAnswer{
								{
									GenerateTime: timestamppb.New(timeNow),
									Answer:       "Tes Answer",
									Feedback:     "LIKE",
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := response.NewDetailAnswerResponse(tt.args.u)
			assert.Equal(t, tt.want.Data.ConversationId, got.Data.ConversationId)
			assert.Equal(t, tt.want.Data.ChatDate, got.Data.ChatDate)
			for k, v := range tt.want.Data.ChatList {
				assert.Equal(t, v.Id, got.Data.ChatList[k].Id)
				assert.Equal(t, v.Question, got.Data.ChatList[k].Question)
				for l, v := range v.Answers {
					assert.Equal(t, v.GenerateTime, got.Data.ChatList[k].Answers[l].GenerateTime)
					assert.Equal(t, v.Answer, got.Data.ChatList[k].Answers[l].Answer)
					assert.Equal(t, v.Feedback, got.Data.ChatList[k].Answers[l].Feedback)
				}
			}
		})
	}
}
