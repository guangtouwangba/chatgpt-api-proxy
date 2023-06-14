package api

import (
	"chatgpt-api-proxy/config"
	"chatgpt-api-proxy/internal/db"
	"net/http"

	"chatgpt-api-proxy/pkg/router"
)

// Handler is the entrypoint for the vercel serverless function.
func Handler(w http.ResponseWriter, r *http.Request) {
	_ = config.NewConfigStore()
	engine := router.NewRouter()
	// init db
	db.SetUpDatabase()
	engine.ServeHTTP(w, r)
}
