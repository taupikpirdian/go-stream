package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/repository"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/internal/config"
	cfg "cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/domain/service"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/exceptions"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/logger"
)

type chatRepository struct {
	cfg         config.ChatConfig
	logger      logger.Logger
	client      *http.Client
	redisClient cfg.RedisClient
}

type ChatRepositoryFactory struct {
	Cfg         config.ChatConfig
	Logger      logger.Logger
	Client      *http.Client
	RedisClient cfg.RedisClient
}

func (d *ChatRepositoryFactory) Create() (repository.ChatRepository, error) {
	if d.Client == nil {
		return nil, exceptions.ErrorRequired("Http client")
	}
	return &chatRepository{
		cfg:         d.Cfg,
		logger:      d.Logger,
		client:      d.Client,
		redisClient: d.RedisClient,
	}, nil
}

func (s *chatRepository) loggerDurationChatbot(ctx context.Context, api string, start time.Time) {
	duration := time.Since(start)
	if duration > 10*time.Second {
		s.logger.ErrorWithContext(fmt.Sprintf("isyana %s API took %f second", api, duration.Seconds()), ctx)
	}
}
