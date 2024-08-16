package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/dto"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/entity"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/internal/repository/service"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/exceptions"
)

func (s *ChatUC) SubmitQuestion(ctx context.Context, dto dto.ChatBotSubmitDto) (*entity.ChatBotSubmit, error) {
	user, err := s.repoUser.GetUsersByID(ctx, dto.UserId)
	if err != nil {
		s.l.Error(err)
		return nil, err
	}
	whitelistUser, err := s.repoWhitelist.GetUsersByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}
	if whitelistUser == nil {
		return nil, exceptions.ErrForbiddenV2
	}

	chatDomain, err := entity.NewSubmitQuestion(dto, user.Name)
	if err != nil {
		return nil, err
	}

	if dto.ConversationId != "" {
		err = s.check24Hours(ctx, dto.UserId, dto.ConversationId, "")
		if err != nil {
			return nil, err
		}
	}

	keyRedis := fmt.Sprintf("chat_genai_%s_%s", chatDomain.RequestId, chatDomain.UserId)
	val, err := s.redisClient.Get(ctx, keyRedis)
	if err != nil {
		if err.Error() != "redis: nil" {
			return nil, err
		}
	}

	response := entity.StartedResponse(chatDomain)
	if val == "" {
		go func() {
			err := s.serviceGenAi.SubmitQuestion(context.Background(), chatDomain)
			if err != nil {
				s.l.Error(err)
			}
		}()
	} else {
		// build from redis
		var output service.WorkflowEventDataWithData
		err = json.Unmarshal([]byte(val), &output)
		if err != nil {
			return nil, err
		}
		response.Status = output.State
		response.ConversationId = output.ConversationID
		response.MessageId = output.MessageID
		response.TaskId = output.TaskID
		response.SubmitDate = output.Data.CreatedAt
		response.SetSubmitDate()

		if response.Status == "ready" {
			if response.Is24Hours() {
				return nil, errors.New("failed to submit question")
			}
			response.Answer = output.Data.Output.Answer
		}
	}

	return response, nil
}

func (s *ChatUC) check24Hours(ctx context.Context, userId, convId, messageId string) error {
	if !s.cfg.SkipCheckConversationId {
		chatDomain, err := entity.NewDetailAnswer(dto.DetailAnswerDto{
			ConversationId: convId,
			UserId:         userId,
		})
		if err != nil {
			return err
		}

		data, err := s.serviceGenAi.DetailAnswer(ctx, chatDomain)
		if err != nil {
			return err
		}

		if data.IsMoreThan24Hour() {
			return errors.New("cant submit/feedback after 24 hour conversation start")
		}

		// check valid messageId when submit feedback
		if messageId != "" {
			found := false
			for _, msg := range data.ChatList {
				if msg.Id == messageId {
					found = true
					break
				}
			}
			if !found {
				return errors.New("invalid message id")
			}
		}
	}
	return nil
}
