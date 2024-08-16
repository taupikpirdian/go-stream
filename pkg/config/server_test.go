package config_test

import (
	"os"
	"strconv"
	"testing"
	"time"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/pkg/config"
)

func TestDefaultValue(t *testing.T) {
	type args struct {
		env          string
		defaultValue string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty env",
			args: args{
				env:          "",
				defaultValue: "localhost",
			},
			want: "localhost",
		},
		{
			name: "have env",
			args: args{
				env:          "DB_HOST",
				defaultValue: "localhost",
			},
			want: "localhost",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.env != "" {
				// Set the environment variable
				os.Setenv(tt.args.env, tt.args.defaultValue)

				// Defer resetting the environment variable
				defer os.Unsetenv(tt.args.env)
			}
			if got := config.DefaultValue(tt.args.env, tt.args.defaultValue); got != tt.want {
				t.Errorf("DefaultValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultValueInt(t *testing.T) {
	type args struct {
		env          string
		defaultValue int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "empty env",
			args: args{
				env:          "",
				defaultValue: 2,
			},
			want: 2,
		},
		{
			name: "have env",
			args: args{
				env:          "APP_TIMEOUT",
				defaultValue: 2,
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.env != "" {
				// Set the environment variable
				os.Setenv(tt.args.env, strconv.Itoa(tt.args.defaultValue))

				// Defer resetting the environment variable
				defer os.Unsetenv(tt.args.env)
			}
			if got := config.DefaultValueInt(tt.args.env, tt.args.defaultValue); got != tt.want {
				t.Errorf("DefaultValueInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultValueDuration(t *testing.T) {
	type args struct {
		env          string
		defaultValue string
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "empty env",
			args: args{
				env:          "",
				defaultValue: "2s",
			},
			want: 2000000000,
		},
		{
			name: "have env",
			args: args{
				env:          "TIME_OUT",
				defaultValue: "2s",
			},
			want: 2000000000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := config.DefaultValueDuration(tt.args.env, tt.args.defaultValue); got != tt.want {
				t.Errorf("DefaultValueDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}
