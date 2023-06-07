package integration_test

import "os"

type IntegrationTestConfig struct {
	APIKey string
	Host   string
}

func InitConfig() *IntegrationTestConfig {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		panic("OPENAI_API_KEY is empty")
	}
	host := os.Getenv("HOST")
	if host == "" {
		panic("HOST is empty")
	}

	return &IntegrationTestConfig{
		APIKey: apiKey,
		Host:   host,
	}
}
