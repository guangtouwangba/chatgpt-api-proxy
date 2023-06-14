package server

import (
	"chatgpt-api-proxy/config"
	"chatgpt-api-proxy/internal/db"
	"chatgpt-api-proxy/pkg/router"

	"github.com/sirupsen/logrus"
)

// RunServer starts the server.
func RunServer() {
	// init config
	store := config.NewConfigStore()
	// init db
	db.SetUpDatabase()

	r := router.NewRouter()
	logrus.Infof("start server at %s", store.Server.Port)
	_ = r.Run(store.GetServerPort())
}
