package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/dto"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/entity"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/repository"
	implementUc "cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/internal/usecase"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/mocks"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/test/testdata"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/exceptions"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestChatUC_ListConversation(t *testing.T) {
	type fields struct {
		serviceGenAi repository.ChatRepository
		userRepo     repository.ChatUsersRepository
		whitelist    repository.WhitelistChatRepository
	}
	type args struct {
		ctx context.Context
		dto dto.ListConversationDto
	}

	var (
		timeNow  = time.Now()
		listData = testdata.NewDomainConversation(timeNow)
	)

	mockUserElig := new(mocks.ChatUsersRepository)
	mockUserElig.On("GetUsersByID", mock.Anything, mock.Anything).
		Return(&entity.User{IsEligibelChatBot: true}, nil)

	mockUserNotElig := new(mocks.ChatUsersRepository)
	mockUserNotElig.On("GetUsersByID", mock.Anything, mock.Anything).
		Return(&entity.User{IsEligibelChatBot: false}, nil)

	mockUserError := new(mocks.ChatUsersRepository)
	mockUserError.On("GetUsersByID", mock.Anything, mock.Anything).
		Return(nil, errors.New("some err"))

	mockGenAi := new(mocks.ChatRepository)
	mockGenAi.On("ListConversation", mock.Anything, mock.Anything).
		Return(nil, errors.New("error hit gen ai"))

	mockGenAi_Success := new(mocks.ChatRepository)
	mockGenAi_Success.On("ListConversation", mock.Anything, mock.Anything).
		Return(listData, nil)

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
		want    *entity.ListConversation
		wantErr bool
	}{
		{
			name: "error - domain",
			fields: fields{
				serviceGenAi: nil,
				userRepo:     mockUserElig,
				whitelist:    mockWhitelist,
			},
			args: args{
				ctx: context.Background(),
				dto: dto.ListConversationDto{
					UserId: "",
				},
			},
			want:    nil,
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
				ctx: context.Background(),
				dto: dto.ListConversationDto{
					UserId: "",
				},
			},
			want:    nil,
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
				ctx: context.Background(),
				dto: dto.ListConversationDto{
					UserId: "",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - get user error",
			fields: fields{
				serviceGenAi: nil,
				userRepo:     mockUserError,
				whitelist:    mockWhitelist,
			},
			args: args{
				ctx: context.Background(),
				dto: dto.ListConversationDto{
					UserId: "",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - usecase",
			fields: fields{
				serviceGenAi: mockGenAi,
				userRepo:     mockUserElig,
				whitelist:    mockWhitelist,
			},
			args: args{
				ctx: context.Background(),
				dto: dto.ListConversationDto{
					UserId: "55101319-0ff4-417f-a4c9-7247e5498d77",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success - usecase",
			fields: fields{
				serviceGenAi: mockGenAi_Success,
				userRepo:     mockUserElig,
				whitelist:    mockWhitelist,
			},
			args: args{
				ctx: context.Background(),
				dto: dto.ListConversationDto{
					UserId: "55101319-0ff4-417f-a4c9-7247e5498d77",
				},
			},
			want:    listData,
			wantErr: false,
		},
		{
			name: "success - 100 data meta",
			fields: fields{
				serviceGenAi: mockGenAi_Success,
				userRepo:     mockUserElig,
				whitelist:    mockWhitelist,
			},
			args: args{
				ctx: context.Background(),
				dto: dto.ListConversationDto{
					UserId: "55101319-0ff4-417f-a4c9-7247e5498d77",
					Limit:  100,
					Page:   1,
				},
			},
			want:    listData,
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

			got, err := d.ListConversation(tt.args.ctx, tt.args.dto)
			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
