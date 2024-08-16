package utils_test

import (
	"crypto/tls"
	"net/http"
	"testing"
	"time"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewHttpClient(t *testing.T) {
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	timeout := 10 * time.Second
	type args struct {
		timeout time.Duration
		proxy   string
	}
	tests := []struct {
		name string
		args args
		want *http.Client
	}{
		{
			name: "success - without proxy",
			args: args{
				timeout: timeout,
				proxy:   "",
			},
			want: &http.Client{Transport: customTransport, Timeout: timeout},
		},
		{
			name: "success - with proxy",
			args: args{
				timeout: timeout,
				proxy:   "localhost:8000",
			},
			want: &http.Client{Transport: customTransport, Timeout: timeout},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.NewHttpClient(tt.args.timeout, tt.args.proxy)
			assert.Equal(t, tt.want.Timeout, got.Timeout)
		})
	}
}

func TestUnixToTime(t *testing.T) {
	timeNow := time.Now()
	timeNowUnix := timeNow.Unix()
	timeNowBangkok := timeNow.Add(time.Hour * 7)

	testCases := []struct {
		name     string
		unix     int64
		expected time.Time
	}{
		{
			name:     "Valid Unix timestamp",
			unix:     int64(timeNowUnix),
			expected: timeNowBangkok,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := utils.UnixToTime(tc.unix)
			assert.Equal(t, tc.expected.Format(time.DateTime), result.Format(time.DateTime))
		})
	}
}
