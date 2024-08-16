package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	cfg "cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/config"
)

const (
	URL                        = "URL"
	TIME_OUT                   = "TIME_OUT"
	CHAT_GENAI_PORT            = "CHAT_GENAI_PORT"
	AUTHORIZATION              = "AUTHORIZATION"
	SKIP_CHECK_CONVERSATION_ID = "SKIP_CHECK_CONVERSATION_ID"
)

type ChatConfig struct {
	Port                    string
	Timeout                 time.Duration
	Url                     string
	Authorization           string
	ProxyUrl                string
	Redis                   *cfg.RedisConfig
	SkipCheckConversationId bool
}

func Chat() ChatConfig {
	err := Validate()
	if err != nil {
		panic(err)
	}

	return ChatConfig{
		Port:                    os.Getenv(CHAT_GENAI_PORT),
		Timeout:                 cfg.DefaultValueDuration(TIME_OUT, "10s"),
		Url:                     os.Getenv(URL),
		Authorization:           os.Getenv(AUTHORIZATION),
		ProxyUrl:                os.Getenv("PROXY_URL"),
		Redis:                   cfg.NewRedisConfig(),
		SkipCheckConversationId: cfg.DefaultValueBool(SKIP_CHECK_CONVERSATION_ID, false),
	}
}

func Validate() error {
	fields := []string{
		URL,
		AUTHORIZATION,
	}

	for _, f := range fields {
		err := Required(f)
		if err != nil {
			return err
		}
	}
	return nil
}

func Required(key string) error {
	if os.Getenv(key) == "" {
		errorMsg := fmt.Sprintf("config %s is required", key)
		return errors.New(errorMsg)
	}
	return nil
}
