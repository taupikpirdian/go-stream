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

type ResponseDetailAnswer struct {
	Limit    int                        `json:"limit"`
	HasMore  bool                       `json:"has_more"`
	Data     []DataResponseDetailAnswer `json:"data"`
	MetaData MetaData                   `json:"metadata"`
}

type DataResponseDetailAnswer struct {
	Id             string `json:"id"`
	ConversationId string `json:"conversation_id"`
	Input          Inputs `json:"inputs"`
	Query          string `json:"query"`
	Answer         string `json:"answer"`
	FeedBack       Rating `json:"feedback"`
	CreatedAt      int64  `json:"created_at"`
	Status         string `json:"status"`
	Error          string `json:"error"`
}

type Rating struct {
	Rating string `json:"rating"`
}

type Inputs struct {
	Name string `json:"name"`
}

type MetaData struct {
	Conversation ConversationMeta `json:"conversation"`
}

type ConversationMeta struct {
	CreatedAt int64 `json:"created_at"`
}

func (s *chatRepository) DetailAnswer(ctx context.Context, reqData *entity.DetailAnswerReq) (*entity.DetailAnswer, error) {
	start := time.Now()
	url := fmt.Sprintf("%s/v1/messages?user=%s&conversation_id=%s&limit=100", s.cfg.Url, reqData.UserId, reqData.ConversationId)
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

	data := &ResponseDetailAnswer{}
	err = json.Unmarshal(resBody, &data)
	if err != nil {
		l.StatusCode = 500
		l.ResponseData = fmt.Sprintf("error UnSerialize: %s", string(resBody))
		return nil, err
	}

	l.StatusCode = 200
	l.RequestData = req.Header
	l.ResponseData = string(resBody)
	s.loggerDurationChatbot(ctx, "detailAnswer", start)

	return MapDetailAnswer(data, reqData.ConversationId), nil
}
