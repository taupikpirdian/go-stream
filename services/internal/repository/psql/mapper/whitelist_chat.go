package mapper

import (
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/entity"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/internal/repository/psql/models"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/common"
)

func ToDomainWhitelistChat(m *models.WhitelistChat) *entity.WhitelistChat {
	id, _ := common.StringToID(m.Id)
	return &entity.WhitelistChat{
		ID:    id,
		Email: m.Email,
	}
}
