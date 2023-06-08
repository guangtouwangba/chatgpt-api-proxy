package config

import (
	"testing"
)

func Test_initConfigs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		want *Config
	}{
		{
			name: "should init config",
			want: &Config{
				Server: serverConfig{
					Port: "8080",
				},
				OpenAI: openAIConfig{
					APIKey: "YOUR_API_KEY",
				},
			},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			NewConfigStore()
			if Store.Server.Port != tc.want.Server.Port {
				t.Errorf("initConfigs() = %v, want %v", Store.GetServerPort(), tc.want.Server.Port)
			}
		})
	}
}
