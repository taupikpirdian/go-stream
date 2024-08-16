package service_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/entity"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/internal/config"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/internal/repository/service"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/test/testdata"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/logger"
	"github.com/stretchr/testify/assert"
)

func Test_chatRepository_DetailAnswer(t *testing.T) {
	type fields struct {
		client *http.Client
	}
	type args struct {
		ctx     context.Context
		reqData *entity.DetailAnswerReq
	}

	var (
		ctx = context.Background()
		cfg = config.ChatConfig{
			Port:          ":8080",
			Timeout:       2,
			Url:           "localhost:9090",
			Authorization: "12121213",
		}
		reqData = &entity.DetailAnswerReq{
			UserId:         "533be47a-9091-4524-afc9-38a9c32447e1",
			ConversationId: "39312c44-c579-427a-b513-33ae492729dc",
		}
		response = `{
			"limit": 2,
			"has_more": true,
			"conversation_created_at": 1714038782,
			"data": [
				{
					"id": "846cfaac-ef01-404b-a8ca-133c4b109307",
					"conversation_id": "39312c44-c579-427a-b513-33ae492729dc",
					"inputs": {
						"name": "Agung Besti"
					},
					"query": "Tes Question",
					"answer": "Tes Answer",
					"message_files": [],
					"feedback": null,
					"retriever_resources": [],
					"created_at": 1714038782,
					"agent_thoughts": [],
					"status": "normal",
					"error": null
				}
			]
		}`
		testdataDomain = testdata.NewDomainDetailAnswer(time.Now())
	)

	fakeClientErrInternal := &http.Client{Transport: testdata.RoundTripFunc(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(strings.NewReader("")),
		}
	})}

	fakeClientResponseError := &http.Client{Transport: testdata.RoundTripFunc(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`this is not json`)),
		}
	})}

	fakeClientOK := &http.Client{Transport: testdata.RoundTripFunc(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(response)),
		}
	})}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.DetailAnswer
		wantErr bool
	}{
		{
			name: "error - process hit",
			fields: fields{
				client: &http.Client{},
			},
			args: args{
				ctx:     ctx,
				reqData: reqData,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - 500",
			fields: fields{
				client: fakeClientErrInternal,
			},
			args: args{
				ctx:     ctx,
				reqData: reqData,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - UnSerialize",
			fields: fields{
				client: fakeClientResponseError,
			},
			args: args{
				ctx:     ctx,
				reqData: reqData,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				client: fakeClientOK,
			},
			args: args{
				ctx:     ctx,
				reqData: reqData,
			},
			want:    testdataDomain,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chatRepoFactory := service.ChatRepositoryFactory{
				Cfg:    cfg,
				Logger: logger.NewFakeApiLogger(),
				Client: tt.fields.client,
			}

			repo, _ := chatRepoFactory.Create()
			data, err := repo.DetailAnswer(tt.args.ctx, tt.args.reqData)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			fmt.Println(data)
			assert.Equal(t, tt.want.ConversationId, data.ConversationId)
			for k, v := range data.ChatList {
				assert.Equal(t, tt.want.ChatList[k].Question, v.Question)
				for kk, vv := range v.Answers {
					assert.Equal(t, tt.want.ChatList[k].Answers[kk].Answer, vv.Answer)
				}
			}
		})
	}
}
