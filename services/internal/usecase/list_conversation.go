package usecase

import (
	"context"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/dto"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/entity"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/exceptions"
)

func (s *ChatUC) ListConversation(ctx context.Context, dto dto.ListConversationDto) (*entity.ListConversation, error) {
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

	chatDomain, err := entity.NewListConversation(dto)
	if err != nil {
		return nil, err
	}

	data, err := s.serviceGenAi.ListConversation(ctx, chatDomain)
	if err != nil {
		return nil, err
	}

	totalData := len(data.Data)
	if totalData > 0 {
		startIndex, endIndex := entity.CalculateDataShow(totalData, int(chatDomain.Limit), int(chatDomain.Page))
		data.Data = data.Data[startIndex:endIndex]
		if chatDomain.Total >= 100 {
			data.Meta.HasMore = false
		}
	}
	return data, nil
}
