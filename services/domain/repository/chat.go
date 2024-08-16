package repository

import (
	"context"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/entity"
)

type ChatRepository interface {
	SubmitQuestion(ctx context.Context, req *entity.ChatBotReq) error
	FeedbackAnswer(ctx context.Context, req *entity.FeedbackReq) error
	StopAnswer(ctx context.Context, req *entity.StopReq) error
	DetailAnswer(ctx context.Context, req *entity.DetailAnswerReq) (*entity.DetailAnswer, error)
	ListConversation(ctx context.Context, req *entity.ListConversationReq) (*entity.ListConversation, error)
}
