package entity_test

import (
	"errors"
	"net/url"
	"strings"
	"testing"
	"time"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/dto"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewSubmitQuestion(t *testing.T) {
	var (
		conversationId = "6b31fdc7-b5ad-44a7-b216-d291fada3a20"
		question       = "test"
		requestId      = "6b31fdc7-b5ad-44a7-b216-d291fada3a21"
	)
	type args struct {
		dto      dto.ChatBotSubmitDto
		userName string
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.ChatBotReq
		wantErr bool
	}{
		{
			name: "error - userId required",
			args: args{
				dto: dto.ChatBotSubmitDto{
					ConversationId: conversationId,
					Question:       question,
					UserId:         "",
					RequestId:      requestId,
				},
				userName: "Yodrick",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - Query required",
			args: args{
				dto: dto.ChatBotSubmitDto{
					ConversationId: conversationId,
					Question:       "",
					UserId:         "Yodrick",
					RequestId:      requestId,
				},
				userName: "Yodrick",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - username required",
			args: args{
				dto: dto.ChatBotSubmitDto{
					ConversationId: conversationId,
					Question:       question,
					UserId:         "Yodrick",
					RequestId:      requestId,
				},
				userName: "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - requestId required",
			args: args{
				dto: dto.ChatBotSubmitDto{
					ConversationId: conversationId,
					Question:       question,
					UserId:         "Yodrick",
					RequestId:      "",
				},
				userName: "Yodrick",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				dto: dto.ChatBotSubmitDto{
					ConversationId: conversationId,
					Question:       question,
					UserId:         "b8aca36e-a4c0-4024-ada8-7ad05e2bf985",
					RequestId:      requestId,
				},
				userName: "Yodrick",
			},
			want: &entity.ChatBotReq{
				UserId:         "b8aca36e-a4c0-4024-ada8-7ad05e2bf985",
				UserName:       "Yodrick",
				Query:          question,
				ConversationId: conversationId,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := entity.NewSubmitQuestion(tt.args.dto, tt.args.userName)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want.UserId, got.UserId)
			assert.Equal(t, tt.want.UserName, got.UserName)
			assert.Equal(t, tt.want.Query, got.Query)
			assert.Equal(t, tt.want.ConversationId, got.ConversationId)
		})
	}
}

func TestNewFeedBackAnswer(t *testing.T) {
	type args struct {
		dto dto.FeedbackAnswerDto
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.FeedbackReq
		wantErr bool
	}{
		{
			name: "error - userId MessageId",
			args: args{
				dto: dto.FeedbackAnswerDto{
					MessageId:      "",
					Feedback:       "like",
					UserId:         "6b31fdc7-b5ad-44a7-b216-d291fada3a21",
					ConversationId: "6b31fdc7-b5ad-44a7-b216-d291fada3a21",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - conversationId",
			args: args{
				dto: dto.FeedbackAnswerDto{
					MessageId: "6b31fdc7-b5ad-44a7-b216-d291fada3a21",
					Feedback:  "LIKE",
					UserId:    "6b31fdc7-b5ad-44a7-b216-d291fada3a21",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - userId Feedback",
			args: args{
				dto: dto.FeedbackAnswerDto{
					MessageId:      "6b31fdc7-b5ad-44a7-b216-d291fada3a21",
					Feedback:       "",
					UserId:         "6b31fdc7-b5ad-44a7-b216-d291fada3a21",
					ConversationId: "6b31fdc7-b5ad-44a7-b216-d291fada3a21",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success - Data Feedback",
			args: args{
				dto: dto.FeedbackAnswerDto{
					MessageId:      "6b31fdc7-b5ad-44a7-b216-d291fada3a21",
					Feedback:       "LIKE",
					UserId:         "6b31fdc7-b5ad-44a7-b216-d291fada3a21",
					ConversationId: "6b31fdc7-b5ad-44a7-b216-d291fada3a21",
				},
			},
			want: &entity.FeedbackReq{
				MessageId:      "6b31fdc7-b5ad-44a7-b216-d291fada3a21",
				Feedback:       "like",
				UserId:         "6b31fdc7-b5ad-44a7-b216-d291fada3a21",
				ConversationId: "6b31fdc7-b5ad-44a7-b216-d291fada3a21",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := entity.NewFeedBackAnswer(tt.args.dto)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want.UserId, got.UserId)
			assert.Equal(t, tt.want.ConversationId, got.ConversationId)
			assert.Equal(t, tt.want.MessageId, got.MessageId)
			assert.Equal(t, tt.want.Feedback, got.Feedback)
		})
	}
}

func TestSetSubmitDate(t *testing.T) {
	epoch := "2024-04-25 10:34:40"
	timeParse, _ := time.Parse(time.DateTime, epoch)
	testCases := []struct {
		name         string
		submitDate   int
		expectedTime time.Time
	}{
		{
			name:         "SubmitDate is a valid epoch",
			submitDate:   int(1714041280),
			expectedTime: timeParse,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := &entity.ChatBotSubmit{
				SubmitDate: tc.submitDate,
			}
			r.SetSubmitDate()
			if !r.SubmitDateTime.Equal(tc.expectedTime) {
				t.Errorf("Expected time %v, got %v", tc.expectedTime, r.SubmitDateTime)
			}
		})
	}
}

func TestStartedResponse(t *testing.T) {
	tests := []struct {
		name     string
		req      *entity.ChatBotReq
		expected *entity.ChatBotSubmit
	}{
		{
			name: "RequestId is not empty",
			req: &entity.ChatBotReq{
				RequestId: "test_request_id",
			},
			expected: &entity.ChatBotSubmit{
				RequestId:      "test_request_id",
				Status:         "pending",
				Feedback:       "NOT_DEFINE",
				SubmitDate:     int(time.Now().Unix()),
				ConversationId: "",
				MessageId:      "",
				TaskId:         "",
				Answer:         "",
			},
		},
		{
			name: "RequestId is empty",
			req:  &entity.ChatBotReq{},
			expected: &entity.ChatBotSubmit{
				Status:         "pending",
				Feedback:       "NOT_DEFINE",
				SubmitDate:     int(time.Now().Unix()),
				ConversationId: "",
				MessageId:      "",
				TaskId:         "",
				Answer:         "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := entity.StartedResponse(tt.req)
			assert.Equal(t, tt.expected.RequestId, res.RequestId)
			assert.Equal(t, tt.expected.Status, res.Status)
			assert.Equal(t, tt.expected.Feedback, res.Feedback)
			assert.NotZero(t, res.SubmitDate)
			assert.Empty(t, res.ConversationId)
			assert.Empty(t, res.MessageId)
			assert.Empty(t, res.TaskId)
			assert.Empty(t, res.Answer)
		})
	}
}

func TestNewStopAnswer(t *testing.T) {
	testCases := []struct {
		name     string
		dto      dto.ChatBotSubmitDto
		taskId   string
		expected *entity.StopReq
		wantErr  bool
	}{
		{
			name: "valid data",
			dto: dto.ChatBotSubmitDto{
				UserId:    "123",
				RequestId: "456",
			},
			taskId: "789",
			expected: &entity.StopReq{
				UserId:    "123",
				RequestId: "456",
				TaskId:    "789",
			},
			wantErr: false,
		},
		{
			name: "invalid data - task id empty",
			dto: dto.ChatBotSubmitDto{
				UserId:    "123",
				RequestId: "456",
			},
			taskId: "",
			expected: &entity.StopReq{
				UserId:    "123",
				RequestId: "456",
				TaskId:    "789",
			},
			wantErr: true,
		},
		{
			name: "invalid data - request id empty",
			dto: dto.ChatBotSubmitDto{
				UserId:    "123",
				RequestId: "",
			},
			taskId: "456",
			expected: &entity.StopReq{
				UserId:    "123",
				RequestId: "",
				TaskId:    "789",
			},
			wantErr: true,
		},
		{
			name:     "invalid data",
			dto:      dto.ChatBotSubmitDto{}, // Empty DTO to trigger validation error
			taskId:   "789",
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := entity.NewStopAnswer(tc.dto, tc.taskId)

			if tc.wantErr && err == nil {
				t.Errorf("Expected error, got nil")
			}

			if !tc.wantErr && err != nil {
				t.Errorf("Expected nil error, got %v", err)
			}

			if !tc.wantErr {
				if result.UserId != tc.expected.UserId {
					t.Errorf("Expected UserId to be %s, got %s", tc.expected.UserId, result.UserId)
				}

				if result.RequestId != tc.expected.RequestId {
					t.Errorf("Expected RequestId to be %s, got %s", tc.expected.RequestId, result.RequestId)
				}

				if result.TaskId != tc.expected.TaskId {
					t.Errorf("Expected TaskId to be %s, got %s", tc.expected.TaskId, result.TaskId)
				}
			}
		})
	}
}

func TestNewDetailAnswer(t *testing.T) {
	testCases := []struct {
		name     string
		dto      dto.DetailAnswerDto
		expected *entity.DetailAnswerReq
		err      error
	}{
		{
			name: "Valid input",
			dto: dto.DetailAnswerDto{
				UserId:         "user1",
				ConversationId: "conversation1",
			},
			expected: &entity.DetailAnswerReq{
				UserId:         "user1",
				ConversationId: "conversation1",
			},
			err: nil,
		},
		{
			name: "Invalid input: UserId is empty",
			dto: dto.DetailAnswerDto{
				UserId:         "",
				ConversationId: "conversation1",
			},
			expected: nil,
			err:      errors.New("UserId is required"),
		},
		{
			name: "Invalid input: ConversationId is empty",
			dto: dto.DetailAnswerDto{
				UserId:         "user1",
				ConversationId: "",
			},
			expected: nil,
			err:      errors.New("ConversationId is required"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := entity.NewDetailAnswer(tc.dto)
			if tc.err != nil {
				assert.NotNil(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, data)

		})
	}
}

func TestNewListConversation(t *testing.T) {
	jakartaTime, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().Local()
	tNow := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, jakartaTime)
	sevenDaysAgo := tNow.AddDate(0, 0, -7)
	encodedDateTime := url.QueryEscape(sevenDaysAgo.Format("2006-01-02 15:04"))
	startDate := strings.ReplaceAll(encodedDateTime, "+", "%20")
	tests := []struct {
		name     string
		dto      dto.ListConversationDto
		expected *entity.ListConversationReq
		err      error
	}{
		{
			name: "Valid input",
			dto: dto.ListConversationDto{
				Page:   1,
				Limit:  10,
				UserId: "user1",
			},
			expected: &entity.ListConversationReq{
				Page:      1,
				Limit:     10,
				UserId:    "user1",
				Total:     10,
				StartDate: startDate,
			},
			err: nil,
		},
		{
			name: "Invalid input: Page is 0",
			dto: dto.ListConversationDto{
				Page:   0,
				Limit:  10,
				UserId: "user1",
			},
			expected: &entity.ListConversationReq{
				Page:      1,
				Limit:     10,
				UserId:    "user1",
				Total:     10,
				StartDate: startDate,
			},
			err: nil,
		},
		{
			name: "Invalid input: Limit is 0",
			dto: dto.ListConversationDto{
				Page:   1,
				Limit:  0,
				UserId: "user1",
			},
			expected: &entity.ListConversationReq{
				Page:      1,
				Limit:     10,
				UserId:    "user1",
				Total:     10,
				StartDate: startDate,
			},
			err: nil,
		},
		{
			name: "Invalid input: UserId is empty",
			dto: dto.ListConversationDto{
				Page:   1,
				Limit:  10,
				UserId: "",
			},
			expected: nil,
			err:      errors.New("UserId is required"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data, err := entity.NewListConversation(tc.dto)

			if tc.err != nil {
				assert.NotNil(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, data)
		})
	}
}

func TestUser_IsEligibelUser(t *testing.T) {
	type fields struct {
		ID                uuid.UUID
		Name              string
		PricePlan         string
		Email             string
		IsEligibelChatBot bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "error",
			fields: fields{
				IsEligibelChatBot: false,
			},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				IsEligibelChatBot: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &entity.User{
				ID:                tt.fields.ID,
				Name:              tt.fields.Name,
				PricePlan:         tt.fields.PricePlan,
				Email:             tt.fields.Email,
				IsEligibelChatBot: tt.fields.IsEligibelChatBot,
			}
			if err := u.IsEligibelUser(); (err != nil) != tt.wantErr {
				t.Errorf("User.IsEligibelUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCalculateDataShow(t *testing.T) {
	type args struct {
		totalData int
		limit     int
		page      int
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 int
	}{
		{
			name: "success",
			args: args{
				totalData: 10,
				limit:     10,
				page:      1,
			},
			want:  0,
			want1: 10,
		},
		{
			name: "success - data to many",
			args: args{
				totalData: 10,
				limit:     10,
				page:      10,
			},
			want:  0,
			want1: 10,
		},
		{
			name: "success - data 100",
			args: args{
				totalData: 100,
				limit:     10,
				page:      10,
			},
			want:  90,
			want1: 100,
		},
		{
			name: "success - data 90 total 100 (limit*page)",
			args: args{
				totalData: 90,
				limit:     10,
				page:      10,
			},
			want:  90,
			want1: 90,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := entity.CalculateDataShow(tt.args.totalData, tt.args.limit, tt.args.page)
			if got != tt.want {
				t.Errorf("CalculateDataShow() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CalculateDataShow() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestChatBotSubmit_Is24Hours(t *testing.T) {
	timeNow := time.Now()
	tests := []struct {
		name           string
		submitDateTime time.Time
		expectedResult bool
	}{
		{
			name:           "Within last 24 hours",
			submitDateTime: timeNow.AddDate(0, 0, -1),
			expectedResult: true,
		},
		{
			name:           "Within last 23 hours",
			submitDateTime: timeNow.Add(time.Hour * -23),
			expectedResult: false,
		},
		{
			name:           "Within last 23 hours 59 minutes 59 seconds",
			submitDateTime: timeNow.Add(time.Hour * -23).Add(time.Minute * -59).Add(time.Second * -59),
			expectedResult: false,
		},
		{
			name:           "In a time",
			submitDateTime: timeNow,
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &entity.ChatBotSubmit{
				SubmitDateTime: tt.submitDateTime,
			}
			if result := r.Is24Hours(); result != tt.expectedResult {
				t.Errorf("Expected %v for %s, but got %v", tt.expectedResult, tt.name, result)
			}
		})
	}
}
