package grpc_server

import (
	"context"
	"time"

	chatPb "cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/pb"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/dto"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/common"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/exceptions"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *chatServiceHandler) StopAnswer(ctx context.Context, req *chatPb.StopAnswerRequest) (*chatPb.MessageResponse, error) {
	l := logger.LogNewResponse{
		RequestStart: time.Now(),
		Routes:       "/v1/T/chat/chatbot/stop",
		RequestData:  req,
	}
	defer l.CreateNewLogV2(true)
	l.TrxId = common.GetMetaDataTransactionId(ctx)

	userMeta, err := common.GetMetaDataUser(ctx)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	dto := dto.ChatBotSubmitDto{
		UserId:    userMeta.Id,
		RequestId: req.GetRequestId(),
	}
	errUc := s.chatUC.StopAnswer(ctx, dto)
	if errUc != nil {
		l.StatusCode = exceptions.MapToGrpcHttpCodeWeb(errUc)
		l.ResponseData = errUc.Error()
		return nil, exceptions.MapToGrpcStatusCodeErr(errUc)
	}

	return &chatPb.MessageResponse{
		Data: &chatPb.MessageData{
			Message: "success",
		},
	}, nil
}
