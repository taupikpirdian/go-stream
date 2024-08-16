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

func (s *chatServiceHandler) FeedbackAnswer(ctx context.Context, req *chatPb.FeedbackAnswerRequest) (*chatPb.MessageResponse, error) {
	l := logger.LogNewResponse{
		RequestStart: time.Now(),
		Routes:       "/v1/T/chat/chatbot/feedback",
		RequestData:  req,
	}
	defer l.CreateNewLogV2(true)
	l.TrxId = common.GetMetaDataTransactionId(ctx)

	userMeta, err := common.GetMetaDataUser(ctx)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	dto := dto.FeedbackAnswerDto{
		MessageId:      req.GetMessageId(),
		Feedback:       req.GetFeedback(),
		UserId:         userMeta.Id,
		ConversationId: req.GetConversationId(),
	}
	errUc := s.chatUC.FeedbackAnswer(ctx, dto)
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
