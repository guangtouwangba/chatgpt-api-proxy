package config

import (
	"chatgpt-api-proxy/pkg/logger"
	"os"

	"github.com/pkg/errors"

	"github.com/spf13/viper"
)

var Store *Config

func NewConfigStore() *Config {
	env := os.Getenv("ENV")
	// default to dev
	if env == "" {
		env = "dev"
	}
	config, err := InitConfig(env)
	if err != nil {
		logger.Panicf("Failed to load config: %v", err)
	}

	return config
}

type Config struct {
	Server   serverConfig
	OpenAI   openAIConfig
	Database DatabaseConfig
}

func (s *Config) GetServerPort() string {
	return ":" + s.Server.Port
}

func (s *Config) GetOpenAIApiKey() string {
	return s.OpenAI.APIKey
}

func (s *Config) GetDatabaseConfig() *DatabaseConfig {
	return &s.Database
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

type DatabaseConfig struct {
	Enabled      bool   `yaml:"enabled"`
	Type         string `yaml:"type"`
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	DatabaseName string `yaml:"databaseName"`
}

func InitConfig(env string) (*Config, error) {
	v, err := loadConfigFile(env)
	if err != nil {
		return nil, err
	}

	serverConfig, err := loadServerConfig(v)
	if err != nil {
		return nil, err
	}

	openAIConfig, err := loadOpenAIConfig(v)
	if err != nil {
		return nil, err
	}

	databaseConfig, err := loadDatabaseConfig(v)
	if err != nil {
		return nil, err
	}

	return &Config{
		Server:   *serverConfig,
		OpenAI:   *openAIConfig,
		Database: *databaseConfig,
	}, nil
}

func loadConfigFile(env string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName("config-" + env)
	v.SetConfigName("config-" + env)
	// default config path
	v.AddConfigPath(".")
	v.AddConfigPath("./conf")
	v.AddConfigPath("../conf")
	v.AddConfigPath("../../conf")
	// test config path
	if env == "test" {
		v.AddConfigPath("../config/testdata")
	}
	// path from environment variable
	configPath := os.Getenv("CONFIG_PATH")
	if configPath != "" {
		v.AddConfigPath(configPath)
	}
	// log output all config path
	logger.Infof("Config path: %v", v.ConfigFileUsed())

	if err := v.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "failed to read config file")
	}
	return v, nil
}

func loadServerConfig(v *viper.Viper) (*serverConfig, error) {
	var c serverConfig
	if err := v.UnmarshalKey("server", &c); err != nil {
		return nil, errors.Wrap(err, "failed to load server config")
	}
	return &c, nil
}

func loadOpenAIConfig(v *viper.Viper) (*openAIConfig, error) {
	var c openAIConfig
	if err := v.UnmarshalKey("openAI", &c); err != nil {
		return nil, errors.Wrap(err, "failed to load openAI config")
	}

	// we prefer use ENV variable for sensitive data
	openAIApiKey := os.Getenv("OPENAI_API_KEY")
	if openAIApiKey != "" {
		logger.Errorf("OpenAI API Key is empty")
	}

	return &c, nil
}

func loadDatabaseConfig(v *viper.Viper) (*DatabaseConfig, error) {
	var c DatabaseConfig
	if err := v.UnmarshalKey("database", &c); err != nil {
		return nil, errors.Wrap(err, "failed to load database config")
	}
	return &c, nil
}
