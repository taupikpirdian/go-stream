package service_test

import (
	"context"
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

func Test_chatRepository_ListConversation(t *testing.T) {
	type fields struct {
		client *http.Client
	}
	type args struct {
		ctx     context.Context
		reqData *entity.ListConversationReq
	}

	var (
		cfg = config.ChatConfig{
			Port:          ":8080",
			Timeout:       2,
			Url:           "localhost:9090",
			Authorization: "12121213",
		}
		ctx = context.Background()
		req = &entity.ListConversationReq{
			Page:      10,
			Limit:     10,
			UserId:    "533be47a-9091-4524-afc9-38a9c32447e1",
			StartDate: "2024-05-14 00:00",
		}
		response = `{
			"limit": 20,
			"has_more": true,
			"data": [
				{
					"id": "39312c44-c579-427a-b513-33ae492729dc",
					"name": "Title 1",
					"inputs": {
						"name": "Agung Besti"
					},
					"status": "normal",
					"introduction": "",
					"created_at": 1714646807
				}
			],
			"metadata": null
		}`
		testdataDomain = testdata.NewDomainConversation(time.Now())
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
		want    *entity.ListConversation
		wantErr bool
	}{
		{
			name: "error - process hit",
			fields: fields{
				client: &http.Client{},
			},
			args: args{
				ctx:     ctx,
				reqData: req,
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
				reqData: req,
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
				reqData: req,
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
				reqData: req,
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
			got, err := repo.ListConversation(tt.args.ctx, tt.args.reqData)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			for i := range got.Data {
				assert.Equal(t, got.Data[i].Id, got.Data[i].Id)
				assert.Equal(t, got.Data[i].Title, got.Data[i].Title)
			}
		})
	}
}
