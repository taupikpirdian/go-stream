package usecase

import (
	"context"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/dto"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/entity"
)

type ChatUseCase interface {
	SubmitQuestion(ctx context.Context, req dto.ChatBotSubmitDto) (*entity.ChatBotSubmit, error)
	FeedbackAnswer(ctx context.Context, req dto.FeedbackAnswerDto) error
	StopAnswer(ctx context.Context, req dto.ChatBotSubmitDto) error
	DetailAnswer(ctx context.Context, req dto.DetailAnswerDto) (*entity.DetailAnswer, error)
	ListConversation(ctx context.Context, req dto.ListConversationDto) (*entity.ListConversation, error)
}
