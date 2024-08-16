package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/exceptions"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/dto"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/entity"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/internal/repository/service"
	"github.com/go-redis/redis/v8"
)

func (s *ChatUC) StopAnswer(ctx context.Context, dto dto.ChatBotSubmitDto) error {
	var (
		taskId = ""
		val    = ""
		err    error
	)

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

	keyRedis := fmt.Sprintf("chat_genai_%s_%s", dto.RequestId, dto.UserId)
	totalCheck := 0
	for {
		totalCheck++
		val, err = s.redisClient.Get(ctx, keyRedis)
		if val != "" {
			break
		} else {
			time.Sleep(1 * time.Second)
		}

		if totalCheck >= 10 {
			break
		}
	}
	if err != nil {
		if err.Error() != redis.Nil.Error() {
			return err
		}
	}

	if val == "" || totalCheck >= 10 {
		errMessage := "redis genai not found"
		s.l.Error(errMessage)
		return errors.New(errMessage)
	} else {
		s.l.Info("redis genai found")
		// build from redis
		var output service.WorkflowEventDataWithData
		err = json.Unmarshal([]byte(val), &output)
		if err != nil {
			return err
		}
		taskId = output.TaskID
	}

	chatDomain, err := entity.NewStopAnswer(dto, taskId)
	if err != nil {
		return err
	}

	err = s.serviceGenAi.StopAnswer(ctx, chatDomain)
	if err != nil {
		return err
	}
	return nil
}
