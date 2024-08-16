package grpc_server

import (
	"context"
	"fmt"
	"time"

	chatPb "cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/pb"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/dto"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/internal/delivery/response"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/common"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/exceptions"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *chatServiceHandler) GetDetailAnswer(ctx context.Context, req *chatPb.GetDetailRequest) (*chatPb.GetDetailResponse, error) {
	l := logger.LogNewResponse{
		RequestStart: time.Now(),
		Routes:       fmt.Sprintf("/v1/T/chat/chatbot/conversation/detail/%s", req.GetConversationId()),
		RequestData:  req,
	}

	defer l.CreateNewLogV2(true)
	l.TrxId = common.GetMetaDataTransactionId(ctx)

	userMeta, err := common.GetMetaDataUser(ctx)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	dto := dto.DetailAnswerDto{
		ConversationId: req.GetConversationId(),
		UserId:         userMeta.Id,
	}
	data, errUc := s.chatUC.DetailAnswer(ctx, dto)
	if errUc != nil {
		l.StatusCode = exceptions.MapToGrpcHttpCodeWeb(errUc)
		l.ResponseData = errUc.Error()
		return nil, exceptions.MapToGrpcStatusCodeErr(errUc)
	}

	return response.NewDetailAnswerResponse(data), nil
}
