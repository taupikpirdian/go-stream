package repository

import (
	"context"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/entity"
)

type ChatUsersRepository interface {
	GetUsersByID(ctx context.Context, userId string) (*entity.User, error)
}
