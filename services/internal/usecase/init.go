package usecase

import (
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/repository"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/usecase"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/internal/config"
	cfg "cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/domain/service"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/exceptions"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/logger"
)

type ChatUCFactory struct {
	L             logger.Logger
	ServiceGenAi  repository.ChatRepository
	RedisClient   cfg.RedisClient
	RepoUser      repository.ChatUsersRepository
	Cfg           config.ChatConfig
	RepoWhitelist repository.WhitelistChatRepository
}

type ChatUC struct {
	l             logger.Logger
	serviceGenAi  repository.ChatRepository
	redisClient   cfg.RedisClient
	repoUser      repository.ChatUsersRepository
	cfg           config.ChatConfig
	repoWhitelist repository.WhitelistChatRepository
}

func (f *ChatUCFactory) Create() (usecase.ChatUseCase, error) {

	if f.L == nil {
		return nil, exceptions.ErrorRequired("logger")
	}

	return &ChatUC{
		l:             f.L,
		serviceGenAi:  f.ServiceGenAi,
		redisClient:   f.RedisClient,
		repoUser:      f.RepoUser,
		cfg:           f.Cfg,
		repoWhitelist: f.RepoWhitelist,
	}, nil

}
