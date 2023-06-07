package config_test

import (
	"testing"

	"chatgpt-api-proxy/config"
)

func Test_initConfigs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		want *config.Config
	}{
		{
			name: "should init config",
			want: &config.Config{
				Server: config.ServerConfig{
					Port: "8080",
				},
				OpenAI: config.OpenAIConfig{
					APIKey: "YOUR_API_KEY",
				},
			},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			config.NewConfigStore()
			if config.Store.Server.Port != tc.want.Server.Port {
				t.Errorf("initConfigs() = %v, want %v", config.Store.GetServerPort(), tc.want.Server.Port)
			}
		})
	}
}
