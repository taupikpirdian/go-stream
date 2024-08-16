package service_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/pkg/utils"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/entity"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/internal/repository/service"
	"github.com/stretchr/testify/assert"
)

func TestManageResponseAndError(t *testing.T) {
	var tests = []struct {
		name           string
		givenAPIName   string
		givenResp      *http.Response
		givenErr       error
		expectedResult []byte
		expectedErr    error
	}{
		{
			name:           "NetworkError",
			givenResp:      nil,
			givenErr:       errors.New("network error"),
			expectedResult: nil,
			expectedErr:    errors.New("network error"),
		},
		{
			name: "NotAuthorized",
			givenResp: &http.Response{
				StatusCode: http.StatusForbidden,
				Body:       httptest.NewRecorder().Result().Body,
			},
			givenErr:       nil,
			expectedResult: []byte{},
			expectedErr:    errors.New("not authorized"),
		},
		{
			name: "Success",
			givenResp: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       httptest.NewRecorder().Result().Body,
			},
			givenErr:       nil,
			expectedResult: []byte{},
			expectedErr:    nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := service.ManageResponseAndError(test.givenResp, test.givenErr)
			assert.ElementsMatch(t, test.expectedResult, result)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func TestMapDetailAnswer(t *testing.T) {
	unix := time.Now().Unix()
	testCases := []struct {
		name           string
		res            *service.ResponseDetailAnswer
		conversationId string
		expected       *entity.DetailAnswer
	}{
		{
			name: "res.Data is not nil",
			res: &service.ResponseDetailAnswer{
				Data: []service.DataResponseDetailAnswer{
					{
						Id:             "533be47a-9091-4524-afc9-38a9c32447e1",
						ConversationId: "533be47a-9091-4524-afc9-38a9c32447e2",
						Input: service.Inputs{
							Name: "test",
						},
						Query:  "Test Question",
						Answer: "Test Answer",
						FeedBack: service.Rating{
							Rating: "LIKE",
						},
						CreatedAt: unix,
					},
				},
			},
			conversationId: "533be47a-9091-4524-afc9-38a9c32447e2",
			expected: &entity.DetailAnswer{
				ConversationId: "533be47a-9091-4524-afc9-38a9c32447e2",
				ChatDate:       utils.UnixToTime(unix),
				ChatList: []*entity.DetailChatList{
					{
						Id:       "533be47a-9091-4524-afc9-38a9c32447e1",
						Question: "Test Question",
						Answers: []*entity.DetailChatListAnswer{
							{
								GenerateTime: utils.UnixToTime(unix),
								Answer:       "Test Answer",
								Feedback:     "LIKE",
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := service.MapDetailAnswer(tc.res, tc.conversationId)
			assert.Equal(t, tc.expected.ConversationId, result.ConversationId)
			for k, v := range tc.expected.ChatList {
				assert.Equal(t, v.Id, result.ChatList[k].Id)
				assert.Equal(t, v.Question, result.ChatList[k].Question)
				assert.Equal(t, v.Answers[0].Answer, result.ChatList[k].Answers[0].Answer)
			}
		})
	}
}
