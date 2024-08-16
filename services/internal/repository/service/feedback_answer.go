package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/entity"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/logger"
)

type BodyFeedbackAnswer struct {
	Rating string `json:"rating"`
	User   string `json:"user"`
}

type ResponseFeedbackAnswer struct {
	Result string `json:"result"`
}

func (s *chatRepository) FeedbackAnswer(ctx context.Context, reqData *entity.FeedbackReq) error {
	start := time.Now()
	url := s.cfg.Url + "/v1/messages/" + reqData.MessageId + "/feedbacks"
	auth := "Bearer " + s.cfg.Authorization

	reqBody := BodyFeedbackAnswer{
		Rating: reqData.Feedback,
		User:   reqData.UserId,
	}
	l := logger.LogNewResponse{
		RequestStart: start,
		Routes:       url,
		RequestData:  reqBody,
	}
	defer l.CreateNewLogV2(false)

	postBody, _ := json.Marshal(reqBody)
	requestBody := bytes.NewBuffer(postBody)

	// Create the HTTP request with the multipart form data
	req, err := http.NewRequest("POST", url, requestBody)
	if err != nil {
		l.StatusCode = 500
		l.ResponseData = fmt.Sprintf("error post: %s", err.Error())
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", auth)

	res, err := s.client.Do(req)
	if err != nil {
		l.StatusCode = 500
		l.ResponseData = fmt.Sprintf("error post: %s", err.Error())
		return err
	}

	resBody, _ := io.ReadAll(res.Body)
	if res.StatusCode != 200 {
		errCode := fmt.Sprintf("error status code: %v", res.StatusCode)
		l.StatusCode = res.StatusCode
		l.ResponseData = string(resBody)
		return errors.New(errCode)
	}

	data := &ResponseFeedbackAnswer{}
	err = json.Unmarshal(resBody, &data)
	if err != nil {
		l.StatusCode = 500
		l.ResponseData = fmt.Sprintf("error UnSerialize: %s", err.Error())
		return err
	}

	errCriteria := data.CriteriaSuccessFeedbackAnswer()
	if errCriteria != nil {
		l.StatusCode = 500
		l.ResponseData = fmt.Sprintf("error Criteria: %s", errCriteria.Error())
		return errCriteria
	}

	l.StatusCode = 200
	l.RequestData = req.Header
	l.ResponseData = string(resBody)
	s.loggerDurationChatbot(ctx, "feedbackAnswer", start)

	return nil
}
