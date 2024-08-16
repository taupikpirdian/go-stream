package grpc_server

import (
	"context"
	"errors"
	"testing"
	"time"

	chatPb "cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/pb"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/usecase"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/mocks"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/test/testdata"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/metadata"
)

func Test_chatServiceHandler_GetDetailAnswer(t *testing.T) {
	type fields struct {
		chatUC usecase.ChatUseCase
	}
	type args struct {
		ctx context.Context
		req *chatPb.GetDetailRequest
	}

	var (
		ctx            = context.Background()
		l              = logger.NewApiLogger(nil)
		timeNow        = time.Now()
		testdataDomain = testdata.NewDomainDetailAnswer(timeNow)
		testdataProtoc = testdata.NewProtoDetailAnswer(timeNow)
	)
	ctxValid := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{"X-USER-ID": "b2b61d92-69f6-4a78-83bb-0aa9d579c6b2", "X-USER-ROLE": "ADMIN"}))

	ucChatBotGenAi_Error := new(mocks.ChatUseCase)
	ucChatBotGenAi_Error.On("DetailAnswer", mock.Anything, mock.Anything).
		Times(1).
		Return(nil, errors.New("error feedback answer"))

	ucChatBotGenAi := new(mocks.ChatUseCase)
	ucChatBotGenAi.On("DetailAnswer", mock.Anything, mock.Anything).
		Times(1).
		Return(testdataDomain, nil)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *chatPb.GetDetailResponse
		wantErr bool
	}{
		{
			name: "error - get user meta",
			fields: fields{
				chatUC: nil,
			},
			args: args{
				ctx: ctx,
				req: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - usecase detail answer",
			fields: fields{
				chatUC: ucChatBotGenAi_Error,
			},
			args: args{
				ctx: ctxValid,
				req: &chatPb.GetDetailRequest{
					ConversationId: "test-conversation-id",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success - usecase detail answer",
			fields: fields{
				chatUC: ucChatBotGenAi,
			},
			args: args{
				ctx: ctxValid,
				req: &chatPb.GetDetailRequest{
					ConversationId: "39312c44-c579-427a-b513-33ae492729dc",
				},
			},
			want:    testdataProtoc,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &chatServiceHandler{
				l:      l,
				chatUC: tt.fields.chatUC,
			}
			got, err := s.GetDetailAnswer(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
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
