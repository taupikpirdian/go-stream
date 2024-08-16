package grpc_server

import (
	"context"
	"time"

	chatPb "cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/pb"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/dto"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/internal/delivery/response"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/common"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/exceptions"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/logger"
)

func (s *chatServiceHandler) SubmitQuestion(ctx context.Context, req *chatPb.SubmitQuestionRequest) (*chatPb.SubmitQuestionResponse, error) {
	l := logger.LogNewResponse{
		RequestStart: time.Now(),
		Routes:       "/v1/T/chat/chatbot/submit",
		RequestData:  req,
	}
	defer l.CreateNewLogV2(true)
	l.TrxId = common.GetMetaDataTransactionId(ctx)

	userMeta, err := common.GetMetaDataUser(ctx)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	dto := dto.ChatBotSubmitDto{
		ConversationId: req.GetConversationId(),
		Question:       req.GetQuestion(),
		UserId:         userMeta.Id,
		RequestId:      req.GetRequestId(),
	}
	res, err := s.chatUC.SubmitQuestion(ctx, dto)
	if err != nil {
		l.StatusCode = exceptions.MapToGrpcHttpCodeWeb(err)
		l.ResponseData = err.Error()
		return nil, exceptions.MapToGrpcStatusCodeErr(err)
	}

	return response.NewSubmitChatResponse(res), nil
}
