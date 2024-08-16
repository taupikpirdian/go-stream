package grpc_server

import (
	"context"
	"errors"
	"testing"

	chatPb "cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/pb"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/usecase"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/mocks"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/test/testdata"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_chatServiceHandler_SubmitQuestion(t *testing.T) {
	var (
		ctx      = context.Background()
		l        = logger.NewApiLogger(nil)
		testdata = testdata.NewAnswerGenAi()
	)
	ctxValid := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{"X-USER-ID": "b2b61d92-69f6-4a78-83bb-0aa9d579c6b2", "X-USER-ROLE": "ADMIN"}))

	ucChatBotGenAi_Error := new(mocks.ChatUseCase)
	ucChatBotGenAi_Error.On("SubmitQuestion", mock.Anything, mock.Anything).
		Times(1).
		Return(nil, errors.New("error submit question"))

	ucChatBotGenAi := new(mocks.ChatUseCase)
	ucChatBotGenAi.On("SubmitQuestion", mock.Anything, mock.Anything).
		Times(1).
		Return(testdata, nil)

	type fields struct {
		l      logger.Logger
		chatUC usecase.ChatUseCase
	}
	type args struct {
		ctx context.Context
		req *chatPb.SubmitQuestionRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *chatPb.SubmitQuestionResponse
		wantErr bool
	}{
		{
			name: "error - get user meta",
			fields: fields{
				l:      nil,
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
			name: "error - usecase submit question",
			fields: fields{
				l:      l,
				chatUC: ucChatBotGenAi_Error,
			},
			args: args{
				ctx: ctxValid,
				req: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success - usecase submit question",
			fields: fields{
				l:      l,
				chatUC: ucChatBotGenAi,
			},
			args: args{
				ctx: ctxValid,
				req: &chatPb.SubmitQuestionRequest{
					ConversationId: "9ec2e82d-51a6-49fd-94ce-b9b24722e286",
					Question:       "test",
				},
			},
			want: &chatPb.SubmitQuestionResponse{
				Data: &chatPb.SubmitQuestionData{
					ConversationId: "9ec2e82d-51a6-49fd-94ce-b9b24722e286",
					SubmitDate:     &timestamppb.Timestamp{},
					Answer:         "Ini contoh jawaban",
					Feedback:       "",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &chatServiceHandler{
				l:      tt.fields.l,
				chatUC: tt.fields.chatUC,
			}
			got, err := s.SubmitQuestion(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want.Data.Answer, got.Data.Answer)
		})
	}
}
