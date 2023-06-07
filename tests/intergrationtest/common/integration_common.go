package integration_test

import "os"

type IntegrationTestConfig struct {
	APIKey string
	Host   string
}

func InitConfig() *IntegrationTestConfig {
	apiKey := os.Getenv("OPENAI_API_KEY")
	host := os.Getenv("HOST")

	return &IntegrationTestConfig{
		APIKey: apiKey,
		Host:   host,
	}
}
