package response_test

import (
	"testing"
	"time"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/test/testdata"
	"github.com/stretchr/testify/assert"

	chatPb "cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/pb"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/entity"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/internal/delivery/response"
)

func TestNewListConversationResponse(t *testing.T) {
	timeNow := time.Now()
	type args struct {
		u *entity.ListConversation
	}
	tests := []struct {
		name string
		args args
		want *chatPb.GetConversationListResponse
	}{
		{
			name: "SUCCESS HASMORE",
			args: args{
				u: testdata.NewDomainConversation(timeNow),
			},
			want: testdata.NewProtoConversation(timeNow),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := response.NewListConversationResponse(tt.args.u)
			for k, v := range tt.want.Data {
				assert.Equal(t, v.Id, got.Data[k].Id)
				assert.Equal(t, v.ChatDate.AsTime(), got.Data[k].ChatDate.AsTime())
				assert.Equal(t, v.Title, got.Data[k].Title)
			}

		})
	}
}
