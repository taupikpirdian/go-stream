package service_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/entity"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/internal/config"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/internal/repository/service"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/test/testdata"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/logger"
	"github.com/stretchr/testify/assert"
)

func Test_chatRepository_FeedbackAnswer(t *testing.T) {
	var (
		ctx = context.Background()
		cfg = config.ChatConfig{
			Port:          ":8080",
			Timeout:       2,
			Url:           "localhost:9090",
			Authorization: "12121213",
		}
		reqData = &entity.FeedbackReq{
			MessageId: "533be47a-9091-4524-afc9-38a9c32447e2",
			Feedback:  "like",
			UserId:    "533be47a-9091-4524-afc9-38a9c32447e1",
		}
		response = `{
			"result": "success"
		}`
		response2 = `{
			"result": "failed"
		}`
	)

	fakeClientErrInternal := &http.Client{Transport: testdata.RoundTripFunc(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(strings.NewReader(response)),
		}
	})}

	fakeClientResponseError := &http.Client{Transport: testdata.RoundTripFunc(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`this is not json`)),
		}
	})}

	fakeClientResponseErrorNotCriteria := &http.Client{Transport: testdata.RoundTripFunc(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(response2)),
		}
	})}

	fakeClientOK := &http.Client{Transport: testdata.RoundTripFunc(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(response)),
		}
	})}

	type fields struct {
		client *http.Client
	}
	type args struct {
		ctx     context.Context
		reqData *entity.FeedbackReq
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
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
			wantErr: true,
		},
		{
			name: "error - not criteria success",
			fields: fields{
				client: fakeClientResponseErrorNotCriteria,
			},
			args: args{
				ctx:     ctx,
				reqData: reqData,
			},
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
			err := repo.FeedbackAnswer(tt.args.ctx, tt.args.reqData)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
