package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChat(t *testing.T) {
	type args struct {
		env []string
		val []string
	}
	tests := []struct {
		name string
		args args
		want ChatConfig
	}{
		{
			name: "success set data config",
			args: args{
				env: []string{"URL", "API_KEY", "SECRET_KEY", "AUTHORIZATION", "REDIS_HOST", "REDIS_PORT", "SKIP_CHECK_CONVERSATION_ID"},
				val: []string{"http://103.175.219.83:8003", "app-CzYznTFqLVuJzYxxtMA", "app-CzYznTFqLVuJzYxxtMB", "app-CzYznTFqLVuJzYxxtMB", "localhost", "6379", "false"},
			},
			want: ChatConfig{
				Port:          ":9090",
				Timeout:       10,
				Url:           "http://103.175.219.83:8003",
				Authorization: "app-CzYznTFqLVuJzYxxtMC",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i, _ := range tt.args.env {
				// Set the environment variable
				os.Setenv(tt.args.env[i], tt.args.val[i])
				// Defer resetting the environment variable
				defer os.Unsetenv(tt.args.env[i])
			}

			assert.Equal(t, tt.want.Url, Chat().Url)
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "error - required env",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
