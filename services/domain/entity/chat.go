package entity

import (
	"errors"
	"net/url"
	"strings"
	"time"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/pkg/utils"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/dto"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/valueobject"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/exceptions"
)

type ChatBotReq struct {
	UserId         string
	UserName       string
	Query          string
	ConversationId string
	RequestId      string
}

type ChatBotSubmit struct {
	ConversationId string
	MessageId      string
	TaskId         string
	Answer         string
	Feedback       string
	SubmitDate     int
	SubmitDateTime time.Time
	RequestId      string
	Status         string
}

type FeedbackReq struct {
	MessageId      string
	Feedback       string
	UserId         string
	ConversationId string
}

type StopReq struct {
	UserId    string
	RequestId string
	TaskId    string
}

type DetailAnswerReq struct {
	ConversationId string
	UserId         string
}

type DetailAnswer struct {
	ConversationId string
	ChatDate       time.Time
	ChatList       []*DetailChatList
}

type DetailChatList struct {
	Id       string
	Question string
	Answers  []*DetailChatListAnswer
}

type DetailChatListAnswer struct {
	GenerateTime time.Time
	Answer       string
	Feedback     string
}

type ListConversation struct {
	Data []*DataConvesation
	Meta MetaData
}

type ListConversationReq struct {
	Page      uint32
	Limit     uint32
	UserId    string
	Total     uint32
	StartDate string
}

type DataConvesation struct {
	Id       string
	Title    string
	ChatDate time.Time
}

type MetaData struct {
	HasMore bool
}

func NewSubmitQuestion(dto dto.ChatBotSubmitDto, userName string) (*ChatBotReq, error) {
	data := &ChatBotReq{
		UserId:         dto.UserId,
		UserName:       userName,
		Query:          dto.Question,
		ConversationId: dto.ConversationId,
		RequestId:      dto.RequestId,
	}
	err := data.validate()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *ChatBotReq) validate() error {
	if r.UserId == "" {
		return exceptions.ErrorRequired("UserId")
	}
	if r.UserName == "" {
		return exceptions.ErrorRequired("UserName")
	}
	if r.Query == "" {
		return exceptions.ErrorRequired("Query")
	}
	if r.RequestId == "" {
		return exceptions.ErrorRequired("RequestId")
	}
	return nil
}

func NewFeedBackAnswer(dto dto.FeedbackAnswerDto) (*FeedbackReq, error) {
	data := &FeedbackReq{
		MessageId:      dto.MessageId,
		Feedback:       valueobject.NewFeedBackFromParam(dto.Feedback).String(),
		UserId:         dto.UserId,
		ConversationId: dto.ConversationId,
	}
	err := data.validate()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *FeedbackReq) validate() error {
	if r.MessageId == "" {
		return exceptions.ErrorRequired("MessageId")
	}
	if r.Feedback == "" {
		return exceptions.ErrorRequired("Feedback")
	}
	if r.ConversationId == "" {
		return exceptions.ErrorRequired("ConversationId")
	}
	return nil
}

func (r *ChatBotSubmit) SetSubmitDate() {
	r.SubmitDateTime = utils.EpochToTime(r.SubmitDate)
}

func StartedResponse(req *ChatBotReq) *ChatBotSubmit {
	currentTime := time.Now()
	return &ChatBotSubmit{
		ConversationId: "",
		MessageId:      "",
		TaskId:         "",
		Answer:         "",
		Feedback:       "NOT_DEFINE",
		SubmitDate:     int(currentTime.Unix()),
		RequestId:      req.RequestId,
		Status:         "pending",
	}
}

func NewStopAnswer(dto dto.ChatBotSubmitDto, taskId string) (*StopReq, error) {
	data := &StopReq{
		UserId:    dto.UserId,
		RequestId: dto.RequestId,
		TaskId:    taskId,
	}
	err := data.validate()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *StopReq) validate() error {
	if r.UserId == "" {
		return exceptions.ErrorRequired("UserId")
	}
	if r.RequestId == "" {
		return exceptions.ErrorRequired("RequestId")
	}
	if r.TaskId == "" {
		return exceptions.ErrorRequired("TaskId")
	}
	return nil
}

func NewDetailAnswer(dto dto.DetailAnswerDto) (*DetailAnswerReq, error) {
	data := &DetailAnswerReq{
		UserId:         dto.UserId,
		ConversationId: dto.ConversationId,
	}
	err := data.validate()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *DetailAnswerReq) validate() error {
	if r.UserId == "" {
		return exceptions.ErrorRequired("UserId")
	}
	if r.ConversationId == "" {
		return exceptions.ErrorRequired("ConversationId")
	}
	return nil
}

func NewListConversation(dto dto.ListConversationDto) (*ListConversationReq, error) {
	data := &ListConversationReq{
		Page:   dto.Page,
		Limit:  dto.Limit,
		UserId: dto.UserId,
	}
	err := data.validate()
	if err != nil {
		return nil, err
	}
	data.Total = data.Page * data.Limit

	jakartaTime, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().Local()
	tNow := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, jakartaTime)
	sevenDaysAgo := tNow.AddDate(0, 0, -7)
	encodedDateTime := url.QueryEscape(sevenDaysAgo.Format("2006-01-02 15:04"))
	data.StartDate = strings.ReplaceAll(encodedDateTime, "+", "%20")
	return data, nil
}

func (r *ListConversationReq) validate() error {
	if r.Page == 0 {
		r.Page = 1
	}
	if r.Limit == 0 {
		r.Limit = 10
	}
	if r.UserId == "" {
		return exceptions.ErrorRequired("UserId")
	}
	return nil
}

func (u *User) IsEligibelUser() error {
	if !u.IsEligibelChatBot {
		return errors.New("not eligible user")
	}
	return nil
}

func CalculateDataShow(totalData, limit, page int) (int, int) {
	startIndex := ((page-1)*limit + 1) - 1
	endIndex := page * limit
	if endIndex > totalData {
		endIndex = totalData
	}
	if startIndex > endIndex {
		startIndex = endIndex - limit
	}
	return startIndex, endIndex
}

func (r *ChatBotSubmit) Is24Hours() bool {
	timeNow := time.Now()
	return timeNow.Unix() >= r.SubmitDateTime.AddDate(0, 0, 1).Unix()
}

func (d *DetailAnswer) IsMoreThan24Hour() bool {
	now := time.Now().Local()
	yesterday := now.AddDate(0, 0, -1)
	if d.ChatDate.Before(yesterday) {
		return true
	}
	return false
}
