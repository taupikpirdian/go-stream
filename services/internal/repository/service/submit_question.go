package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/pkg/exception"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/entity"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/logger"
)

type Data struct {
	EventData WorkflowEventData `json:"data"`
}

type WorkflowEventData struct {
	Event          string `json:"event"`
	ConversationID string `json:"conversation_id"`
	MessageID      string `json:"message_id"`
	CreatedAt      int64  `json:"created_at"`
	TaskID         string `json:"task_id"`
	WorkflowRunID  string `json:"workflow_run_id"`
	State          string
}

type WorkflowEventDataWithData struct {
	Event          string `json:"event"`
	ConversationID string `json:"conversation_id"`
	MessageID      string `json:"message_id"`
	CreatedAt      int64  `json:"created_at"`
	TaskID         string `json:"task_id"`
	WorkflowRunID  string `json:"workflow_run_id"`
	State          string
	Data           DataWorkFlow `json:"data"`
}

type DataWorkFlow struct {
	Id             string    `json:"id"`
	WorkflowId     string    `json:"workflow_id"`
	SequenceNumber int       `json:"sequence_number"`
	Status         string    `json:"status"`
	Output         Output    `json:"outputs"`
	Error          string    `json:"error"`
	ElapsedTime    float32   `json:"elapsed_time"`
	TotalTokens    int       `json:"total_tokens"`
	TotalSteps     int       `json:"total_steps"`
	CreatedBy      CreatedBy `json:"created_by"`
	CreatedAt      int       `json:"created_at"`
	FinishedAt     int       `json:"finished_at"`
}

type Output struct {
	Answer string `json:"answer"`
}

type CreatedBy struct {
	Id   string `json:"id"`
	User string `json:"user"`
}

type BodyChatSubmit struct {
	Inputs         InputName `json:"inputs"`
	Query          string    `json:"query"`
	ConversationId string    `json:"conversation_id"`
	User           string    `json:"user"`
	ResponseMode   string    `json:"response_mode"`
}

type InputName struct {
	Name string `json:"name"`
}

func (s *chatRepository) SubmitQuestion(ctx context.Context, reqData *entity.ChatBotReq) error {
	start := time.Now()
	url := s.cfg.Url + "/v1/chat-messages"
	auth := "Bearer " + s.cfg.Authorization
	keyRedis := fmt.Sprintf("chat_genai_%s_%s", reqData.RequestId, reqData.UserId)
	var isSuccess = false
	var workflows string

	reqBody := BodyChatSubmit{
		Inputs: InputName{
			Name: reqData.UserName,
		},
		Query:          reqData.Query,
		ConversationId: reqData.ConversationId,
		User:           reqData.UserId,
		ResponseMode:   "streaming",
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
	req, _ := http.NewRequest("POST", url, requestBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", auth)

	res, err := s.client.Do(req)
	if err != nil {
		l.StatusCode = 500
		l.ResponseData = fmt.Sprintf("error post: %s", err.Error())

		err = s.setRedis(ctx, keyRedis, "", "failed", reqData)
		if err != nil {
			return err
		}
		return err
	}

	if res.StatusCode != 200 {
		l.StatusCode = res.StatusCode
		errCode := fmt.Sprintf("error status code: %v", res.StatusCode)
		l.ResponseData = errCode

		err = s.setRedis(ctx, keyRedis, "", "failed", reqData)
		if err != nil {
			return err
		}
		return errors.New(errCode)
	}
	reader := bufio.NewReader(res.Body)
	for {
		var d WorkflowEventData
		line, err := reader.ReadString('\n')
		if err != nil {
			s.logger.Error(fmt.Sprintf("error reading string: %s", err.Error()))
			break
		}
		dataString := strings.TrimPrefix(line, "data: ")
		isValidString := ValidString(dataString)
		if !isValidString {
			continue
		}

		err = json.Unmarshal([]byte(dataString), &d)
		if err != nil {
			s.logger.Warn(fmt.Sprintf("error unmarshal: %s", err.Error()))
			continue
		}
		workflows += d.Event + "|"
		if d.Event == "workflow_started" {
			d.State = "pending"
			outputString, err := json.Marshal(d)
			if err != nil {
				break
			}
			err = s.redisClient.Set(ctx, keyRedis, outputString, 15*time.Minute)
			if err != nil {
				s.logger.Error(err)
				break
			}
		} else if d.Event == "workflow_finished" {
			isSuccess = true
			err := s.setRedis(ctx, keyRedis, dataString, "success", reqData)
			if err != nil {
				s.logger.Error(err)
				break
			}
		} else if d.Event == "message_end" {
			fmt.Println("message_end...")
			break
		}
	}

	l.StatusCode = 200
	l.RequestData = req.Header
	s.loggerDurationChatbot(ctx, "submitQuestionChatBot", start)
	if !isSuccess {
		s.logger.Info(workflows)
		s.logger.Warn(time.Since(start))
		err = s.setRedis(ctx, keyRedis, "", "failed", reqData)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *chatRepository) setRedis(ctx context.Context, key string, value string, typeSet string, reqData *entity.ChatBotReq) error {
	var err error
	var outputString []byte
	if typeSet == "failed" {
		output := WorkflowEventData{
			Event:          "",
			ConversationID: reqData.ConversationId,
			MessageID:      "",
			CreatedAt:      time.Now().Unix(),
			TaskID:         "",
			WorkflowRunID:  "",
			State:          "failed",
		}
		outputString, err = json.Marshal(output)
		if err != nil {
			s.logger.Error(err)
			return err
		}
	} else if typeSet == "success" {
		var output WorkflowEventDataWithData
		err = json.Unmarshal([]byte(value), &output)
		if err != nil {
			s.logger.Error(err)
			return err
		}

		output.State = "ready"
		outputString, err = json.Marshal(output)
		if err != nil {
			s.logger.Error(err)
			return err
		}
	}

	err = s.redisClient.Set(ctx, key, outputString, 15*time.Minute)
	if err != nil {
		s.logger.Error(err)
	}
	return err
}

func ValidString(responseString string) bool {
	if len(responseString) == 1 {
		return false
	}
	if responseString == exception.ErrGenAiEventPing.Error() {
		return false
	}
	return true
}
