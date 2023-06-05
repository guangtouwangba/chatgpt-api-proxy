package config

import (
	"chatgpt-api-proxy/config"
	"testing"
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
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.NewConfigStore()
			if config.Store.Server.Port != tt.want.Server.Port {
				t.Errorf("initConfigs() = %v, want %v", config.Store.GetServerPort(), tt.want.Server.Port)
			}
		})
	}
}
