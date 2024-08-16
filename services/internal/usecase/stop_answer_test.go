package usecase_test

import (
	"context"
	"errors"
	"testing"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/entity"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/dto"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/repository"
	implementUc "cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/internal/usecase"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/mocks"
	cfg "cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/domain/service"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/exceptions"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/logger"
	mockPkg "cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestChatUC_StopAnswer(t *testing.T) {
	type fields struct {
		serviceGenAi  repository.ChatRepository
		redisClient   cfg.RedisClient
		userRepo      repository.ChatUsersRepository
		whitelistRepo repository.WhitelistChatRepository
	}
	type args struct {
		ctx context.Context
		dto dto.ChatBotSubmitDto
	}

	var (
		ctx = context.Background()
		val = "{\"event\":\"workflow_finished\",\"conversation_id\":\"1325fa7b-09c9-4363-880b-d1d9d90f7f52\",\"message_id\":\"191c7871-5f3e-4014-9d1b-9cbfa2c9491c\",\"created_at\":1714381439,\"task_id\":\"3b206cae-ce26-4509-b12d-47422ad5b6b4\",\"workflow_run_id\":\"a9d23e34-3fe2-4152-9abd-cb880bed793f\",\"State\":\"success\",\"data\":{\"id\":\"a9d23e34-3fe2-4152-9abd-cb880bed793f\",\"workflow_id\":\"5846d192-0f16-4f30-906b-7cee63bd46f8\",\"sequence_number\":256,\"status\":\"succeeded\",\"outputs\":{\"answer\":\"Untuk pengaturan display logic di tSurvey, Anda dapat mengikuti langkah-langkah berikut:\\n\\n1. Masuk ke halaman survei yang ingin Anda edit.\\n2. Pilih pertanyaan yang ingin Anda atur logika tampilannya.\\n3. Pada bagian bawah pertanyaan, klik \\\"Add Destination Page Logic\\\" untuk menambahkan logika tampilan.\\n4. Anda dapat mengatur logika tampilan berdasarkan jawaban yang dipilih oleh responden sebelumnya.\\n5. Setelah mengatur logika tampilan, klik \\\"Save\\\" untuk menyimpan perubahan.\\n\\nJika Anda mengalami kendala atau membutuhkan bantuan lebih lanjut, jangan ragu untuk menghubungi kami. Kami siap membantu Anda.\"},\"error\":\"\",\"elapsed_time\":6.7047825,\"total_tokens\":2479,\"total_steps\":6,\"created_by\":{\"id\":\"a02da62d-8f6d-463d-9a29-c047afefe9bd\",\"user\":\"e520c3ee-b339-4871-b0e6-74235e907f6b\"},\"created_at\":1714381439,\"finished_at\":1714381445}}"
	)

	redis_Error := new(mockPkg.RedisClient)
	redis_Error.On("Get", mock.Anything, mock.Anything).
		Return("", errors.New("error"))

	redis_ErrorNotFound := new(mockPkg.RedisClient)
	redis_ErrorNotFound.On("Get", mock.Anything, mock.Anything).
		Return("", nil)

	redis_ErrorNotJson := new(mockPkg.RedisClient)
	redis_ErrorNotJson.On("Get", mock.Anything, mock.Anything).
		Return("have string", nil)

	redis := new(mockPkg.RedisClient)
	redis.On("Get", mock.Anything, mock.Anything).
		Return(val, nil)

	mockGenAi := new(mocks.ChatRepository)
	mockGenAi.On("StopAnswer", mock.Anything, mock.Anything).
		Times(1).
		Return(errors.New("error hit gen ai"))

	mockGenAi_Success := new(mocks.ChatRepository)
	mockGenAi_Success.On("StopAnswer", mock.Anything, mock.Anything).
		Times(1).
		Return(nil)

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
			name: "error - repo user",
			fields: fields{
				serviceGenAi: nil,
				redisClient:  redis_Error,
				userRepo:     mockUserElig_Error,
			},
			args: args{
				ctx: ctx,
				dto: dto.ChatBotSubmitDto{
					RequestId: "6b31fdc7-b5ad-44a7-b216-d291fada3a20",
					UserId:    "6b31fdc7-b5ad-44a7-b216-d291fada3a21",
				},
			},
			wantErr: true,
		},
		{
			name: "error - repo whitelist",
			fields: fields{
				serviceGenAi:  nil,
				redisClient:   redis_Error,
				userRepo:      mockUserElig,
				whitelistRepo: mockWhitelistError,
			},
			args: args{
				ctx: ctx,
				dto: dto.ChatBotSubmitDto{
					RequestId: "6b31fdc7-b5ad-44a7-b216-d291fada3a20",
					UserId:    "6b31fdc7-b5ad-44a7-b216-d291fada3a21",
				},
			},
			wantErr: true,
		},
		{
			name: "error - user whitelist not eligible",
			fields: fields{
				serviceGenAi:  nil,
				redisClient:   redis_Error,
				userRepo:      mockUserElig,
				whitelistRepo: mockWhitelistNotEligible,
			},
			args: args{
				ctx: ctx,
				dto: dto.ChatBotSubmitDto{
					RequestId: "6b31fdc7-b5ad-44a7-b216-d291fada3a20",
					UserId:    "6b31fdc7-b5ad-44a7-b216-d291fada3a21",
				},
			},
			wantErr: true,
		},
		{
			name: "error - redis",
			fields: fields{
				serviceGenAi:  nil,
				redisClient:   redis_Error,
				userRepo:      mockUserElig,
				whitelistRepo: mockWhitelist,
			},
			args: args{
				ctx: ctx,
				dto: dto.ChatBotSubmitDto{
					RequestId: "6b31fdc7-b5ad-44a7-b216-d291fada3a20",
					UserId:    "6b31fdc7-b5ad-44a7-b216-d291fada3a21",
				},
			},
			wantErr: true,
		},
		{
			name: "error - redis not found",
			fields: fields{
				serviceGenAi:  nil,
				redisClient:   redis_ErrorNotFound,
				userRepo:      mockUserElig,
				whitelistRepo: mockWhitelist,
			},
			args: args{
				ctx: ctx,
				dto: dto.ChatBotSubmitDto{
					RequestId: "6b31fdc7-b5ad-44a7-b216-d291fada3a20",
					UserId:    "6b31fdc7-b5ad-44a7-b216-d291fada3a21",
				},
			},
			wantErr: true,
		},
		{
			name: "error - unmarshal",
			fields: fields{
				serviceGenAi:  nil,
				redisClient:   redis_ErrorNotJson,
				userRepo:      mockUserElig,
				whitelistRepo: mockWhitelist,
			},
			args: args{
				ctx: ctx,
				dto: dto.ChatBotSubmitDto{
					RequestId: "6b31fdc7-b5ad-44a7-b216-d291fada3a20",
					UserId:    "6b31fdc7-b5ad-44a7-b216-d291fada3a21",
				},
			},
			wantErr: true,
		},
		{
			name: "error - domain",
			fields: fields{
				serviceGenAi:  nil,
				redisClient:   redis,
				userRepo:      mockUserElig,
				whitelistRepo: mockWhitelist,
			},
			args: args{
				ctx: ctx,
				dto: dto.ChatBotSubmitDto{
					RequestId: "6b31fdc7-b5ad-44a7-b216-d291fada3a20",
					UserId:    "",
				},
			},
			wantErr: true,
		},
		{
			name: "error - usecase",
			fields: fields{
				serviceGenAi:  mockGenAi,
				redisClient:   redis,
				userRepo:      mockUserElig,
				whitelistRepo: mockWhitelist,
			},
			args: args{
				ctx: ctx,
				dto: dto.ChatBotSubmitDto{
					RequestId: "6b31fdc7-b5ad-44a7-b216-d291fada3a20",
					UserId:    "6b31fdc7-b5ad-44a7-b216-d291fada3a21",
				},
			},
			wantErr: true,
		},
		{
			name: "success - usecase",
			fields: fields{
				serviceGenAi:  mockGenAi_Success,
				redisClient:   redis,
				userRepo:      mockUserElig,
				whitelistRepo: mockWhitelist,
			},
			args: args{
				ctx: ctx,
				dto: dto.ChatBotSubmitDto{
					RequestId: "6b31fdc7-b5ad-44a7-b216-d291fada3a20",
					UserId:    "6b31fdc7-b5ad-44a7-b216-d291fada3a21",
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
				RedisClient:   tt.fields.redisClient,
				RepoUser:      tt.fields.userRepo,
				RepoWhitelist: tt.fields.whitelistRepo,
			}
			d, err := factory.Create()
			assert.Nil(t, err)

			err = d.StopAnswer(tt.args.ctx, tt.args.dto)
			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
		})
	}
}
