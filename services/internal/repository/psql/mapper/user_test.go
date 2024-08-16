package mapper_test

import (
	"testing"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/entity"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/internal/repository/psql/mapper"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/internal/repository/psql/models"
	"github.com/stretchr/testify/assert"
)

func TestToDomainUser(t *testing.T) {
	type args struct {
		m *models.User
	}
	tests := []struct {
		name string
		args args
		want *entity.User
	}{
		{
			name: "",
			args: args{
				m: &models.User{
					Name:                 "Yor",
					Email:                "yor@gmail.com",
					PricePlan:            "tsel",
					IsEligibleForChatbot: false,
				},
			},
			want: &entity.User{
				Name:              "Yor",
				PricePlan:         "tsel",
				Email:             "yor@gmail.com",
				IsEligibelChatBot: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapper.ToDomainUser(tt.args.m)
			assert.Equal(t, tt.want.Name, got.Name)
			assert.Equal(t, tt.want.Email, got.Email)
			assert.Equal(t, tt.want.PricePlan, got.PricePlan)
			assert.Equal(t, tt.want.IsEligibelChatBot, got.IsEligibelChatBot)
		})
	}
}
