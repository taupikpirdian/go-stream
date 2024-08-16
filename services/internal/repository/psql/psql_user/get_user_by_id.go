package psql_user

import (
	"context"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/entity"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/internal/repository/psql/mapper"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/internal/repository/psql/models"
)

func (r *userRepository) GetUsersByID(ctx context.Context, userId string) (*entity.User, error) {
	var mdl *models.User
	tx := r.db.WithContext(ctx).
		Where("id = ?", userId).
		First(&mdl)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return mapper.ToDomainUser(mdl), nil
}
