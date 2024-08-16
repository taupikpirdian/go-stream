package usecase

import (
	"context"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/dto"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/entity"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/exceptions"
)

func (s *ChatUC) DetailAnswer(ctx context.Context, dto dto.DetailAnswerDto) (*entity.DetailAnswer, error) {
	user, err := s.repoUser.GetUsersByID(ctx, dto.UserId)
	if err != nil {
		return nil, err
	}

	whitelistUser, err := s.repoWhitelist.GetUsersByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}
	if whitelistUser == nil {
		return nil, exceptions.ErrForbiddenV2
	}

	chatDomain, err := entity.NewDetailAnswer(dto)
	if err != nil {
		return nil, err
	}

	data, err := s.serviceGenAi.DetailAnswer(ctx, chatDomain)
	if err != nil {
		return nil, err
	}
	return data, nil
}
