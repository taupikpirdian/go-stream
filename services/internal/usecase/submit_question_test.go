package usecase_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/internal/repository/service"

	cfg "cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/domain/service"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/exceptions"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/dto"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/entity"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/repository"
	implementUc "cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/internal/usecase"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/mocks"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/test/testdata"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/logger"
	mockPkg "cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestChatUC_SubmitQuestion(t *testing.T) {
	var (
		ctx          = context.Background()
		dataDtoError = dto.ChatBotSubmitDto{
			ConversationId: "6b31fdc7-b5ad-44a7-b216-d291fada3a20",
			Question:       "test",
			RequestId:      "6b31fdc7-b5ad-44a7-b216-d291fada3a22",
		}
		dataDto = dto.ChatBotSubmitDto{
			ConversationId: "6b31fdc7-b5ad-44a7-b216-d291fada3a20",
			Question:       "test",
			UserId:         "6b31fdc7-b5ad-44a7-b216-d291fada3a21",
			RequestId:      "6b31fdc7-b5ad-44a7-b216-d291fada3a22",
		}
		responseGenAi   = testdata.NewAnswerGenAi()
		val             = "{\"event\":\"workflow_finished\",\"conversation_id\":\"9ec2e82d-51a6-49fd-94ce-b9b24722e287\",\"message_id\":\"191c7871-5f3e-4014-9d1b-9cbfa2c9491c\",\"created_at\":1714381439,\"task_id\":\"3b206cae-ce26-4509-b12d-47422ad5b6b4\",\"workflow_run_id\":\"a9d23e34-3fe2-4152-9abd-cb880bed793f\",\"State\":\"ready\",\"data\":{\"id\":\"a9d23e34-3fe2-4152-9abd-cb880bed793f\",\"workflow_id\":\"5846d192-0f16-4f30-906b-7cee63bd46f8\",\"sequence_number\":256,\"status\":\"succeeded\",\"outputs\":{\"answer\":\"Untuk pengaturan display logic di tSurvey, Anda dapat mengikuti langkah-langkah berikut:\\n\\n1. Masuk ke halaman survei yang ingin Anda edit.\\n2. Pilih pertanyaan yang ingin Anda atur logika tampilannya.\\n3. Pada bagian bawah pertanyaan, klik \\\"Add Destination Page Logic\\\" untuk menambahkan logika tampilan.\\n4. Anda dapat mengatur logika tampilan berdasarkan jawaban yang dipilih oleh responden sebelumnya.\\n5. Setelah mengatur logika tampilan, klik \\\"Save\\\" untuk menyimpan perubahan.\\n\\nJika Anda mengalami kendala atau membutuhkan bantuan lebih lanjut, jangan ragu untuk menghubungi kami. Kami siap membantu Anda.\"},\"error\":\"\",\"elapsed_time\":6.7047825,\"total_tokens\":2479,\"total_steps\":6,\"created_by\":{\"id\":\"a02da62d-8f6d-463d-9a29-c047afefe9bd\",\"user\":\"e520c3ee-b339-4871-b0e6-74235e907f6b\"},\"created_at\":1714381439,\"finished_at\":1714381445}}"
		user            = testdata.NewUser(true)
		userNotEligibel = testdata.NewUser(false)
		valWithTimeNow  = &service.WorkflowEventDataWithData{
			Event:          "workflow_finished",
			ConversationID: "9ec2e82d-51a6-49fd-94ce-b9b24722e287",
			MessageID:      "191c7871-5f3e-4014-9d1b-9cbfa2c9491c",
			CreatedAt:      0,
			TaskID:         "3b206cae-ce26-4509-b12d-47422ad5b6b4",
			WorkflowRunID:  "a9d23e34-3fe2-4152-9abd-cb880bed793f",
			State:          "ready",
			Data: service.DataWorkFlow{
				Id:             "a9d23e34-3fe2-4152-9abd-cb880bed793f",
				WorkflowId:     "5846d192-0f16-4f30-906b-7cee63bd46f8",
				SequenceNumber: 256,
				Status:         "succeeded",
				Output: service.Output{
					Answer: "Ini contoh jawaban",
				},
				Error:       "",
				ElapsedTime: 6.7047825,
				TotalTokens: 2479,
				TotalSteps:  6,
				CreatedBy: service.CreatedBy{
					Id:   "a02da62d-8f6d-463d-9a29-c047afefe9bd",
					User: "e520c3ee-b339-4871-b0e6-74235e907f6b",
				},
				CreatedAt:  int(time.Now().Unix()),
				FinishedAt: int(time.Now().Unix()),
			},
		}
		dataTest          = testdata.NewDomainDetailAnswer(time.Now())
		dataTestYesterday = testdata.NewDomainDetailAnswer(time.Now().Add(-25 * time.Hour))
	)

	mockGenAi := new(mocks.ChatRepository)
	mockGenAi.On("DetailAnswer", mock.Anything, mock.Anything).
		Times(5).
		Return(dataTest, nil)
	mockGenAi.On("SubmitQuestion", mock.Anything, mock.Anything).
		Times(1).
		Return(nil, errors.New("error hit gen ai"))

	mockGenAiErrRedis := new(mocks.ChatRepository)
	mockGenAiErrRedis.On("DetailAnswer", mock.Anything, mock.Anything).
		Times(1).
		Return(dataTest, nil)

	mockGenAi_Success := new(mocks.ChatRepository)
	mockGenAi_Success.On("SubmitQuestion", mock.Anything, mock.Anything).
		Times(1).
		Return(responseGenAi, nil)
	mockGenAi_Success.On("DetailAnswer", mock.Anything, mock.Anything).
		Times(1).
		Return(dataTest, nil)

	mockGenAi_ErrGetDetail := new(mocks.ChatRepository)
	mockGenAi_ErrGetDetail.On("DetailAnswer", mock.Anything, mock.Anything).
		Times(1).
		Return(nil, errors.New("some err"))

	mockGenAi_moreThan1Days := new(mocks.ChatRepository)
	mockGenAi_moreThan1Days.On("DetailAnswer", mock.Anything, mock.Anything).
		Times(1).
		Return(dataTestYesterday, nil)

	redis_Success := new(mockPkg.RedisClient)
	redis_Success.On("Get", mock.Anything, mock.Anything).
		Return(val, nil)

	outputString, _ := json.Marshal(valWithTimeNow)
	redis_Success_Valid := new(mockPkg.RedisClient)
	redis_Success_Valid.On("Get", mock.Anything, mock.Anything).
		Return(string(outputString), nil)

	redis_ErrorNotObject := new(mockPkg.RedisClient)
	redis_ErrorNotObject.On("Get", mock.Anything, mock.Anything).
		Return("not a object", nil)

	redis_Error := new(mockPkg.RedisClient)
	redis_Error.On("Get", mock.Anything, mock.Anything).
		Return("", errors.New("error"))

	repoUser_Error := new(mocks.ChatUsersRepository)
	repoUser_Error.On("GetUsersByID", mock.Anything, mock.Anything).
		Return(nil, errors.New("error"))

	repoUser_NotEligibel := new(mocks.ChatUsersRepository)
	repoUser_NotEligibel.On("GetUsersByID", mock.Anything, mock.Anything).
		Return(userNotEligibel, nil)

	repoUser := new(mocks.ChatUsersRepository)
	repoUser.On("GetUsersByID", mock.Anything, mock.Anything).
		Return(user, nil)

	mockWhitelist := new(mocks.WhitelistChatRepository)
	mockWhitelist.On("GetUsersByEmail", mock.Anything, mock.Anything).
		Return(&entity.WhitelistChat{Email: "test@gmail.com"}, nil)

	mockWhitelistError := new(mocks.WhitelistChatRepository)
	mockWhitelistError.On("GetUsersByEmail", mock.Anything, mock.Anything).
		Return(nil, exceptions.ErrInternalServerError)

	mockWhitelistNotEligible := new(mocks.WhitelistChatRepository)
	mockWhitelistNotEligible.On("GetUsersByEmail", mock.Anything, mock.Anything).
		Return(nil, nil)

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
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.ChatBotSubmit
		wantErr bool
	}{
		{
			name: "error - find user",
			fields: fields{
				serviceGenAi: nil,
				userRepo:     repoUser_Error,
			},
			args: args{
				ctx: ctx,
				dto: dataDtoError,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - repo whitelist",
			fields: fields{
				serviceGenAi:  nil,
				userRepo:      repoUser,
				whitelistRepo: mockWhitelistError,
			},
			args: args{
				ctx: ctx,
				dto: dataDtoError,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - user whitelist not eligible",
			fields: fields{
				serviceGenAi:  nil,
				userRepo:      repoUser,
				whitelistRepo: mockWhitelistNotEligible,
			},
			args: args{
				ctx: ctx,
				dto: dataDtoError,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - domain",
			fields: fields{
				serviceGenAi:  nil,
				userRepo:      repoUser,
				whitelistRepo: mockWhitelist,
			},
			args: args{
				ctx: ctx,
				dto: dataDtoError,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - redis",
			fields: fields{
				serviceGenAi:  mockGenAiErrRedis,
				redisClient:   redis_Error,
				userRepo:      repoUser,
				whitelistRepo: mockWhitelist,
			},
			args: args{
				ctx: ctx,
				dto: dataDto,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - Unmarshal",
			fields: fields{
				serviceGenAi:  mockGenAi,
				redisClient:   redis_ErrorNotObject,
				userRepo:      repoUser,
				whitelistRepo: mockWhitelist,
			},
			args: args{
				ctx: ctx,
				dto: dataDto,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - conversation expired, only can submit before 24 hours",
			fields: fields{
				serviceGenAi:  mockGenAi,
				redisClient:   redis_Success,
				userRepo:      repoUser,
				whitelistRepo: mockWhitelist,
			},
			args: args{
				ctx: ctx,
				dto: dataDto,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - cant submit after 24 hours conversation started",
			fields: fields{
				serviceGenAi:  mockGenAi_moreThan1Days,
				redisClient:   redis_Success_Valid,
				userRepo:      repoUser,
				whitelistRepo: mockWhitelist,
			},
			args: args{
				ctx: ctx,
				dto: dataDto,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - failed get conversation detail",
			fields: fields{
				serviceGenAi:  mockGenAi_ErrGetDetail,
				redisClient:   redis_Success_Valid,
				userRepo:      repoUser,
				whitelistRepo: mockWhitelist,
			},
			args: args{
				ctx: ctx,
				dto: dataDto,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success - ready service gen ai",
			fields: fields{
				serviceGenAi:  mockGenAi_Success,
				redisClient:   redis_Success_Valid,
				userRepo:      repoUser,
				whitelistRepo: mockWhitelist,
			},
			args: args{
				ctx: ctx,
				dto: dataDto,
			},
			want:    responseGenAi,
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

			got, err := d.SubmitQuestion(tt.args.ctx, tt.args.dto)
			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, tt.want.RequestId, got.RequestId)
			assert.Equal(t, tt.want.ConversationId, got.ConversationId)
			assert.Equal(t, tt.want.Answer, got.Answer)
			assert.Equal(t, tt.want.Feedback, got.Feedback)
		})
	}
}
