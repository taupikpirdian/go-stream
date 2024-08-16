package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/entity"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/test/testdata"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/dto"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/repository"
	implementUc "cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/internal/usecase"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/mocks"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/exceptions"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestChatUC_FeedbackAnswer(t *testing.T) {
	var (
		ctx               = context.Background()
		dataTest          = testdata.NewDomainDetailAnswer(time.Now())
		dataTestYesterday = testdata.NewDomainDetailAnswer(time.Now().Add(-25 * time.Hour))
	)

	mockGenAi := new(mocks.ChatRepository)
	mockGenAi.On("FeedbackAnswer", mock.Anything, mock.Anything).
		Times(1).
		Return(errors.New("error hit gen ai"))
	mockGenAi.On("DetailAnswer", mock.Anything, mock.Anything).
		Times(1).
		Return(dataTest, nil)

	mockGenAi_Success := new(mocks.ChatRepository)
	mockGenAi_Success.On("FeedbackAnswer", mock.Anything, mock.Anything).
		Times(1).
		Return(nil)
	mockGenAi_Success.On("DetailAnswer", mock.Anything, mock.Anything).
		Times(1).
		Return(dataTest, nil)

	mockGenAi_Error := new(mocks.ChatRepository)
	mockGenAi_Error.On("FeedbackAnswer", mock.Anything, mock.Anything).
		Times(1).
		Return(exceptions.ErrInternalServerError)
	mockGenAi_Error.On("DetailAnswer", mock.Anything, mock.Anything).
		Times(1).
		Return(dataTest, nil)

	mockGenAiYesterday := new(mocks.ChatRepository)
	mockGenAiYesterday.On("DetailAnswer", mock.Anything, mock.Anything).
		Times(1).
		Return(dataTestYesterday, nil)

	type fields struct {
		serviceGenAi repository.ChatRepository
		userRepo     repository.ChatUsersRepository
		whitelist    repository.WhitelistChatRepository
	}
	type args struct {
		ctx context.Context
		dto dto.FeedbackAnswerDto
	}

	mockUserElig := new(mocks.ChatUsersRepository)
	mockUserElig.On("GetUsersByID", mock.Anything, mock.Anything).
		Return(&entity.User{IsEligibelChatBot: true}, nil)

	mockUserElig_False := new(mocks.ChatUsersRepository)
	mockUserElig_False.On("GetUsersByID", mock.Anything, mock.Anything).
		Return(&entity.User{IsEligibelChatBot: false}, nil)

	mockUserElig_Error := new(mocks.ChatUsersRepository)
	mockUserElig_Error.On("GetUsersByID", mock.Anything, mock.Anything).
		Return(nil, errors.New("error hit gen ai"))

	mockWhitelist := new(mocks.WhitelistChatRepository)
	mockWhitelist.On("GetUsersByEmail", mock.Anything, mock.Anything).
		Return(&entity.WhitelistChat{Email: "test@gmail.com"}, nil)

	mockWhitelistError := new(mocks.WhitelistChatRepository)
	mockWhitelistError.On("GetUsersByEmail", mock.Anything, mock.Anything).
		Return(nil, exceptions.ErrInternalServerError)

	mockWhitelistNotEligible := new(mocks.WhitelistChatRepository)
	mockWhitelistNotEligible.On("GetUsersByEmail", mock.Anything, mock.Anything).
		Return(nil, nil)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "error - repository failure",
			fields: fields{
				serviceGenAi: nil,
				userRepo:     mockUserElig_Error,
			},
			args: args{
				ctx: nil,
				dto: dto.FeedbackAnswerDto{
					MessageId: "",
					Feedback:  "LIKE",
					UserId:    "",
				},
			},
			wantErr: true,
		},
		{
			name: "error - repo whitelist",
			fields: fields{
				serviceGenAi: nil,
				userRepo:     mockUserElig,
				whitelist:    mockWhitelistError,
			},
			args: args{
				ctx: nil,
				dto: dto.FeedbackAnswerDto{
					MessageId: "",
					Feedback:  "LIKE",
					UserId:    "",
				},
			},
			wantErr: true,
		},
		{
			name: "error - user whitelist not eligible",
			fields: fields{
				serviceGenAi: nil,
				userRepo:     mockUserElig,
				whitelist:    mockWhitelistNotEligible,
			},
			args: args{
				ctx: nil,
				dto: dto.FeedbackAnswerDto{
					MessageId: "",
					Feedback:  "LIKE",
					UserId:    "",
				},
			},
			wantErr: true,
		},
		{
			name: "error - domain",
			fields: fields{
				serviceGenAi: nil,
				userRepo:     mockUserElig,
				whitelist:    mockWhitelist,
			},
			args: args{
				ctx: nil,
				dto: dto.FeedbackAnswerDto{
					MessageId: "",
					Feedback:  "LIKE",
					UserId:    "",
				},
			},
			wantErr: true,
		},
		{
			name: "error - hit gen ai",
			fields: fields{
				serviceGenAi: mockGenAi,
				userRepo:     mockUserElig,
				whitelist:    mockWhitelist,
			},
			args: args{
				ctx: ctx,
				dto: dto.FeedbackAnswerDto{
					MessageId:      "55101319-0ff4-417f-a4c9-7247e5498d77",
					Feedback:       "LIKE",
					UserId:         "55101319-0ff4-417f-a4c9-7247e5498d71",
					ConversationId: "55101319-0ff4-417f-a4c9-7247e5498d71",
				},
			},
			wantErr: true,
		},
		{
			name: "error - more than 24 hour cant send feedback",
			fields: fields{
				serviceGenAi: mockGenAiYesterday,
				userRepo:     mockUserElig,
				whitelist:    mockWhitelist,
			},
			args: args{
				ctx: ctx,
				dto: dto.FeedbackAnswerDto{
					MessageId:      "55101319-0ff4-417f-a4c9-7247e5498d77",
					Feedback:       "LIKE",
					UserId:         "55101319-0ff4-417f-a4c9-7247e5498d71",
					ConversationId: "55101319-0ff4-417f-a4c9-7247e5498d71",
				},
			},
			wantErr: true,
		},
		{
			name: "error - feedback answer",
			fields: fields{
				serviceGenAi: mockGenAi_Error,
				userRepo:     mockUserElig,
				whitelist:    mockWhitelist,
			},
			args: args{
				ctx: ctx,
				dto: dto.FeedbackAnswerDto{
					MessageId:      "39312c44-c579-427a-b513-33ae492729d1",
					Feedback:       "LIKE",
					UserId:         "55101319-0ff4-417f-a4c9-7247e5498d71",
					ConversationId: "55101319-0ff4-417f-a4c9-7247e5498d71",
				},
			},
			wantErr: true,
		},
		{
			name: "success - hit gen ai",
			fields: fields{
				serviceGenAi: mockGenAi_Success,
				userRepo:     mockUserElig,
				whitelist:    mockWhitelist,
			},
			args: args{
				ctx: ctx,
				dto: dto.FeedbackAnswerDto{
					MessageId:      "39312c44-c579-427a-b513-33ae492729d1",
					Feedback:       "LIKE",
					UserId:         "55101319-0ff4-417f-a4c9-7247e5498d71",
					ConversationId: "55101319-0ff4-417f-a4c9-7247e5498d71",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := implementUc.ChatUCFactory{
				L:             logger.NewFakeApiLogger(),
				ServiceGenAi:  tt.fields.serviceGenAi,
				RepoUser:      tt.fields.userRepo,
				RepoWhitelist: tt.fields.whitelist,
			}
			d, err := factory.Create()
			assert.Nil(t, err)

			err = d.FeedbackAnswer(tt.args.ctx, tt.args.dto)
			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
		})
	}
}
