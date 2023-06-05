package config

import (
	"log"
	"os"

)

var Store *Config

func NewConfigStore() *Config {
	initConfigs()

	return Store
}

type Config struct {
	Server ServerConfig
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

func (s *Config) GetServerPort() string {
	return ":" + s.Server.Port
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
