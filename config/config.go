package config

import (
	"log"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var Store *Config

func NewConfigStore() *Config {
	initConfigs()

	return Store
}

type Config struct {
	Server ServerConfig
	OpenAI OpenAIConfig
}

func (s *Config) GetServerPort() string {
	return ":" + s.Server.Port
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type OpenAIConfig struct {
	Type     string `yaml:"type"`
	APIKey   string `yaml:"apiKey"`
	UserName string `yaml:"userName"`
	Password string `yaml:"password"`
	Model    string `yaml:"model"`
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
		logrus.Info("fatal error config file: default \n", err)
		os.Exit(1)
	}

	port := viper.GetString("server.port")

	logrus.Info("port: ", port)

	if err := viper.Unmarshal(&Store); err != nil {
		log.Fatal(err)
	}
}
