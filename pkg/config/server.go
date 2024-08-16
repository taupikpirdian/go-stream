package config

import (
	"os"
	"strconv"
	"time"
)

func DefaultValue(env string, defaultValue string) string {
	val := os.Getenv(env)
	if val == "" {
		return defaultValue
	}
	return val
}

func DefaultValueInt(env string, defaultValue int) int {
	val := ConvertInt(env)
	if val == 0 {
		return defaultValue
	}
	return val
}

func ConvertInt(env string) int {
	v, _ := strconv.Atoi(os.Getenv(env))
	return v
}

func DefaultValueDuration(env string, defaultValue string) time.Duration {
	value := os.Getenv(env)
	if value == "" {
		value = defaultValue
	}

	d, err := time.ParseDuration(value)
	if err != nil {
		panic(err)
	}
	return d
}
