package api

import (
	"net/http"

	"chatgpt-api-proxy/pkg/router"
)

// Handler is the entrypoint for the vercel serverless function.
func Handler(w http.ResponseWriter, r *http.Request) {
	engine := router.NewRouter()

	engine.ServeHTTP(w, r)
}
