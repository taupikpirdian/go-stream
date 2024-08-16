package service

import (
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

type ResponseConversation struct {
	Limit   int                        `json:"limit"`
	HasMore bool                       `json:"has_more"`
	Data    []DataResponseConversation `json:"data"`
}

type DataResponseConversation struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Input        Inputs `json:"inputs"`
	Status       string `json:"status"`
	Introduction string `json:"introduction"`
	CreatedAt    int64  `json:"created_at"`
}

func (s *chatRepository) ListConversation(ctx context.Context, reqData *entity.ListConversationReq) (*entity.ListConversation, error) {
	start := time.Now()
	url := fmt.Sprintf("%s/v1/conversations?user=%s&limit=%v&start=%v", s.cfg.Url, reqData.UserId, reqData.Total, reqData.StartDate)
	auth := "Bearer " + s.cfg.Authorization
	l := logger.LogNewResponse{
		RequestStart: start,
		Routes:       url,
		RequestData:  reqData,
	}
	defer l.CreateNewLogV2(false)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", auth)

	res, err := s.client.Do(req)
	if err != nil {
		l.StatusCode = 500
		l.ResponseData = fmt.Sprintf("error post: %s", err.Error())
		return nil, err
	}

	resBody, _ := io.ReadAll(res.Body)
	if res.StatusCode != 200 {
		errCode := fmt.Sprintf("error status code: %v", res.StatusCode)
		l.StatusCode = res.StatusCode
		l.ResponseData = string(resBody)
		return nil, errors.New(errCode)
	}

	data := &ResponseConversation{}
	err = json.Unmarshal(resBody, &data)
	if err != nil {
		l.StatusCode = 500
		l.ResponseData = fmt.Sprintf("error UnSerialize: %s", err.Error())
		return nil, err
	}

	l.StatusCode = 200
	l.RequestData = req.Header
	l.ResponseData = string(resBody)
	s.loggerDurationChatbot(ctx, "listConversation", start)

	return MapConversation(data), nil
}
