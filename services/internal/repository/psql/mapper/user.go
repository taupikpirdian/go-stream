package mapper

import (
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/entity"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/internal/repository/psql/models"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/common"
)

func ToDomainUser(m *models.User) *entity.User {
	id, _ := common.StringToID(m.Id)
	return &entity.User{
		ID:                id,
		Name:              m.Name,
		PricePlan:         m.PricePlan,
		Email:             m.Email,
		IsEligibelChatBot: m.IsEligibleForChatbot,
	}
}
