package repository

import (
	"context"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/entity"
)

type WhitelistChatRepository interface {
	GetUsersByEmail(ctx context.Context, email string) (*entity.WhitelistChat, error)
}
