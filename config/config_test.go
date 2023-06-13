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

func TestLoadConfig(t *testing.T) {
	// 调用待测试的函数
	config, err := InitConfig("test")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// server port should be 8080
	if config.Server.Port != "8080" {
		t.Fatalf("Server port should be 8080")
	}

	// if database is enabled, then database host should be localhost
	if config.Database.Enabled && config.Database.Host != "localhost" {
		t.Fatalf("Database host should be localhost")
	}

	// if database is enabled, port should be 5432
	if config.Database.Enabled && config.Database.Port != "5432" {
		t.Fatalf("Database port should be 5432")
	}
}
