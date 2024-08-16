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

func Test_chatServiceHandler_GetConversationList(t *testing.T) {
	type fields struct {
		chatUC usecase.ChatUseCase
	}
	type args struct {
		ctx context.Context
		req *chatPb.GetConversationListRequest
	}

	var (
		ctx            = context.Background()
		l              = logger.NewApiLogger(nil)
		timeNow        = time.Now()
		testdataDomain = testdata.NewDomainConversation(timeNow)
		testdataProtoc = testdata.NewProtoConversation(timeNow)
	)
	ctxValid := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{"X-USER-ID": "b2b61d92-69f6-4a78-83bb-0aa9d579c6b2", "X-USER-ROLE": "ADMIN"}))

	ucChatBotGenAi_Error := new(mocks.ChatUseCase)
	ucChatBotGenAi_Error.On("ListConversation", mock.Anything, mock.Anything).
		Times(1).
		Return(nil, errors.New("error feedback answer"))

	ucChatBotGenAi := new(mocks.ChatUseCase)
	ucChatBotGenAi.On("ListConversation", mock.Anything, mock.Anything).
		Times(1).
		Return(testdataDomain, nil)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *chatPb.GetConversationListResponse
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
			name: "error - usecase list conversation",
			fields: fields{
				chatUC: ucChatBotGenAi_Error,
			},
			args: args{
				ctx: ctxValid,
				req: &chatPb.GetConversationListRequest{
					Limit: 10,
					Page:  1,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success - usecase list conversation",
			fields: fields{
				chatUC: ucChatBotGenAi,
			},
			args: args{
				ctx: ctxValid,
				req: &chatPb.GetConversationListRequest{
					Limit: 10,
					Page:  1,
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
			got, err := s.GetConversationList(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			for k, v := range got.Data {
				assert.Equal(t, tt.want.Data[k].Id, v.Id)
				assert.Equal(t, tt.want.Data[k].Title, v.Title)
				assert.Equal(t, tt.want.Data[k].ChatDate, v.ChatDate)
			}
		})
	}
}
