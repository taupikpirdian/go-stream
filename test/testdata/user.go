package testdata

import (
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/entity"
)

func NewUser(eligibel bool) *entity.User {
	return &entity.User{
		Name:              "Yor",
		PricePlan:         "tsel",
		Email:             "yor@gmail.com",
		IsEligibelChatBot: eligibel,
	}
}
