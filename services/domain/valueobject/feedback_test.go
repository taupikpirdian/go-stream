package valueobject_test

import (
	"reflect"
	"testing"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/valueobject"
)

func TestNewFeedBackFromParam(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want valueobject.Feedback
	}{
		{
			name: "NOT DEFINE",
			args: args{
				s: "",
			},
			want: valueobject.NewFeedBack(valueobject.NOT_DEFINE),
		},
		{
			name: "LIKE",
			args: args{
				s: "LIKE",
			},
			want: valueobject.NewFeedBack(valueobject.LIKE),
		},
		{
			name: "DISLIKE",
			args: args{
				s: "DISLIKE",
			},
			want: valueobject.NewFeedBack(valueobject.DISLIKE),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := valueobject.NewFeedBackFromParam(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFeedBackFromParam() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFeedback_String(t *testing.T) {
	tests := []struct {
		name     string
		feedback valueobject.Feedback
		want     string
	}{
		{
			name:     "LIKE",
			feedback: valueobject.NewFeedBackFromParam("LIKE"),
			want:     "like",
		},
		{
			name:     "DISLIKE",
			feedback: valueobject.NewFeedBackFromParam("DISLIKE"),
			want:     "dislike",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.feedback.String(); got != tt.want {
				t.Errorf("Feedback.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
