package grpc_server

import (
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/usecase"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/logger"
	"google.golang.org/grpc"
)

type Server struct {
	GRPCServer *grpc.Server
	L          logger.Logger
	ChatUC     usecase.ChatUseCase
}

func (s *Server) CreateHandlerChatService() (*chatServiceHandler, error) {
	handler := &chatServiceHandler{
		s:      s.GRPCServer,
		l:      s.L,
		chatUC: s.ChatUC,
	}
	return handler, nil
}
