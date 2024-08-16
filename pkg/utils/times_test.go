package utils_test

import (
	"testing"
	"time"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestEpochToTime(t *testing.T) {
	epoch := "2024-04-25 17:34:40"
	timeParse, _ := time.Parse(time.DateTime, epoch)

	type args struct {
		epoch int
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "",
			args: args{
				epoch: 1714041280,
			},
			want: timeParse,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.EpochToTime(tt.args.epoch)
			assert.Equal(t, tt.want.Format(time.DateOnly), got.Format(time.DateOnly))
		})
	}
}
