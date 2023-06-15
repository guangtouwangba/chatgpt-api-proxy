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
	var config *Config
	config, err := InitConfigFromConfigFile(env)
	if err != nil {
		logger.Errorf("Failed to load config: %v", err)
		// if we can't load the config, we try to load the config from env
		config = InitConfigFromEnv()
	}
	Store = config
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

func InitConfigFromEnv() *Config {
	viper.AutomaticEnv()

	// server config
	serverPort := viper.GetString("SERVER_PORT")
	if serverPort == "" {
		// default to 8080
		// we do not panic here because for some platforms like vercel, we do not need to set SERVER_PORT
		logger.Errorf("SERVER_PORT is not set")
	}

	// openAI config
	openAIApiKey := viper.GetString("OPENAI_API_KEY")
	if openAIApiKey == "" {
		// we do not panic here because we allow user to pass OPENAI_API_KEY via http header
		logger.Errorf("OPENAI_API_KEY is not set")
	}

	// database config
	dbConfig := DatabaseConfig{}
	databaseEnabled := viper.GetBool("DATABASE_ENABLED")
	if databaseEnabled {
		dbConfig = loadDBConfig(databaseEnabled)
	}

	return &Config{
		Server:   serverConfig{Port: serverPort},
		OpenAI:   openAIConfig{APIKey: openAIApiKey},
		Database: dbConfig,
	}
}

func loadDBConfig(databaseEnabled bool) DatabaseConfig {
	databaseType := viper.GetString("DATABASE_TYPE")
	if databaseType == "" {
		logger.Errorf("DATABASE_TYPE is not set")
	}
	databaseHost := viper.GetString("DATABASE_HOST")
	if databaseHost == "" {
		logger.Panicf("DATABASE_HOST is not set")
	}
	databasePort := viper.GetString("DATABASE_PORT")
	if databasePort == "" {
		logger.Panicf("DATABASE_PORT is not set")
	}
	databaseUser := viper.GetString("DATABASE_USER")
	if databaseUser == "" {
		logger.Panicf("DATABASE_USER is not set")
	}
	databasePassword := viper.GetString("DATABASE_PASSWORD")
	if databasePassword == "" {
		logger.Panicf("DATABASE_PASSWORD is not set")
	}
	databaseName := viper.GetString("DATABASE_NAME")
	if databaseName == "" {
		logger.Panicf("DATABASE_NAME is not set")
	}
	return DatabaseConfig{
		Enabled:      databaseEnabled,
		Type:         databaseType,
		Host:         databaseHost,
		Port:         databasePort,
		User:         databaseUser,
		Password:     databasePassword,
		DatabaseName: databaseName,
	}
}

// InitConfigFromConfigFile loads the config from a config file
// if you are using vercel, you should not use this function to load config.
func InitConfigFromConfigFile(env string) (*Config, error) {
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
	if openAIApiKey == "" {
		logger.Errorf("OpenAI API Key is empty")
	}

	return &c, nil
}

func loadDatabaseConfig(v *viper.Viper) (*DatabaseConfig, error) {
	var c DatabaseConfig
	var password string
	if err := v.UnmarshalKey("database", &c); err != nil {
		return nil, errors.Wrap(err, "failed to load database config")
	}

	if env := os.Getenv("ENV"); env == "prod" {
		password = os.Getenv("DB_PASSWORD")
		c.Password = password
	}
	return &c, nil
}
