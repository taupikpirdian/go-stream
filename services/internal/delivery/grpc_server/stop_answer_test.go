package grpc_server

import (
	"context"
	"errors"
	"testing"

	chatPb "cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/pb"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/usecase"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/mocks"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/metadata"
)

func Test_chatServiceHandler_StopAnswer(t *testing.T) {
	var (
		ctx = context.Background()
		l   = logger.NewApiLogger(nil)
		// testdata = testdata.NewAnswerGenAi()
	)
	ctxValid := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{"X-USER-ID": "b2b61d92-69f6-4a78-83bb-0aa9d579c6b2", "X-USER-ROLE": "ADMIN"}))

	ucChatBotGenAi_Error := new(mocks.ChatUseCase)
	ucChatBotGenAi_Error.On("StopAnswer", mock.Anything, mock.Anything).
		Times(1).
		Return(errors.New("error feedback answer"))

	ucChatBotGenAi := new(mocks.ChatUseCase)
	ucChatBotGenAi.On("StopAnswer", mock.Anything, mock.Anything).
		Times(1).
		Return(nil)

	type fields struct {
		chatUC usecase.ChatUseCase
	}
	type args struct {
		ctx context.Context
		req *chatPb.StopAnswerRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *chatPb.MessageResponse
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
			name: "error - usecase stop answer",
			fields: fields{
				chatUC: ucChatBotGenAi_Error,
			},
			args: args{
				ctx: ctxValid,
				req: &chatPb.StopAnswerRequest{
					RequestId: "799c209a-4805-4c77-b43d-b4361c2b6af4",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success - usecase stop answer",
			fields: fields{
				chatUC: ucChatBotGenAi,
			},
			args: args{
				ctx: ctxValid,
				req: &chatPb.StopAnswerRequest{
					RequestId: "799c209a-4805-4c77-b43d-b4361c2b6af4",
				},
			},
			want: &chatPb.MessageResponse{
				Data: &chatPb.MessageData{
					Message: "success",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &chatServiceHandler{
				l:      l,
				chatUC: tt.fields.chatUC,
			}
			got, err := s.StopAnswer(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want.Data.Message, got.Data.Message)
		})
	}
}
