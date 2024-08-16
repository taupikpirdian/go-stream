package grpc_server

import (
	chatPb "cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/pb"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/usecase"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/logger"
	"google.golang.org/grpc"
)

type chatServiceHandler struct {
	s      *grpc.Server
	l      logger.Logger
	chatUC usecase.ChatUseCase
}

func (h *chatServiceHandler) Handle() {
	chatPb.RegisterChatServiceServer(h.s, h)
}
