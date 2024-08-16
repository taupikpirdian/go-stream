package usecase

import (
	"context"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/exceptions"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/dto"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/entity"
)

func (s *ChatUC) FeedbackAnswer(ctx context.Context, dto dto.FeedbackAnswerDto) error {
	user, err := s.repoUser.GetUsersByID(ctx, dto.UserId)
	if err != nil {
		return err
	}
	whitelistUser, err := s.repoWhitelist.GetUsersByEmail(ctx, user.Email)
	if err != nil {
		return err
	}
	if whitelistUser == nil {
		return exceptions.ErrForbiddenV2
	}

	chatDomain, err := entity.NewFeedBackAnswer(dto)
	if err != nil {
		return err
	}

	err = s.check24Hours(ctx, dto.UserId, dto.ConversationId, dto.MessageId)
	if err != nil {
		return err
	}

	err = s.serviceGenAi.FeedbackAnswer(ctx, chatDomain)
	if err != nil {
		return err
	}
	return nil
}
