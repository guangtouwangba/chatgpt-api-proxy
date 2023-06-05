package server

import (
	"chatgpt-api-proxy/config"
	"chatgpt-api-proxy/internal/router"
	"github.com/sirupsen/logrus"
)

func RunServer() {
	// init config
	store := config.NewConfigStore()

	r := router.NewRouter()
	logrus.Infof("start server at %s", store.Server.Port)
	_ = r.Run(store.GetServerPort())
}
