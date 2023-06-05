package api

import (
	"net/http"

	"chatgpt-api-proxy/config"
	"chatgpt-api-proxy/pkg/router"
	"github.com/sirupsen/logrus"
)

// Handler is the entrypoint for the vercel serverless function.
func Handler(w http.ResponseWriter, r *http.Request) {
	// init config
	store := config.NewConfigStore()

	engine := router.NewRouter()
	logrus.Infof("start server at %s", store.Server.Port)

	engine.ServeHTTP(w, r)
}
