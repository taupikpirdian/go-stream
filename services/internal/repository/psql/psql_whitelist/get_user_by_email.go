package psql_whitelist

import (
	"context"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/entity"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/internal/repository/psql/mapper"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/internal/repository/psql/models"
)

func (r *whiteListRepository) GetUsersByEmail(ctx context.Context, email string) (*entity.WhitelistChat, error) {
	var mdl *models.WhitelistChat
	tx := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&mdl)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return mapper.ToDomainWhitelistChat(mdl), nil
}
