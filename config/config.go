package config

import (
	"chatgpt-api-proxy/pkg/logger"
	"os"

	"github.com/spf13/viper"
)

var Store *Config

func NewConfigStore() *Config {
	initConfigs()

	return Store
}

type Config struct {
	Server serverConfig
	OpenAI openAIConfig
}

func (s *Config) GetServerPort() string {
	return ":" + s.Server.Port
}

func (s *Config) GetOpenAIApiKey() string {
	return s.OpenAI.APIKey
}

type serverConfig struct {
	Port  string `yaml:"port"`
	Proxy string `yaml:"proxy"`
}

type openAIConfig struct {
	Type     string `yaml:"type"`
	APIKey   string `yaml:"apiKey"`
	UserName string `yaml:"userName"`
	Password string `yaml:"password"`
}

func initConfigs() {
	viper.SetConfigName("default") // config file name without extension
	viper.SetConfigType("yaml")
	// TODO: refactor config path
	viper.AddConfigPath(".")
	viper.AddConfigPath("./conf")
	viper.AddConfigPath("../conf")
	viper.AddConfigPath("../../conf")
	viper.AutomaticEnv() // read value ENV variable

	err := viper.ReadInConfig()
	if err != nil {
		logger.Errorf("Fatal error config file: %s \n", err)
		os.Exit(1)
	}

	port := viper.GetString("server.port")

	logger.Infof("Server port: %s", port)

	if err := viper.Unmarshal(&Store); err != nil {
		logger.Errorf("unable to decode into struct, %v", err)
	}
	// we prefer use ENV variable for sensitive data
	openAIApiKey := os.Getenv("OPENAI_API_KEY")
	if openAIApiKey != "" {
		Store.OpenAI.APIKey = openAIApiKey
	}
	logger.Infof("OpenAI API Key: %s", Store.OpenAI.APIKey)
}
