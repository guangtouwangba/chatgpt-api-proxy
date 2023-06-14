package api

import (
	"chatgpt-api-proxy/pkg/server"
	"net/http"
)

// Handler is the entrypoint for the vercel serverless function.
func Handler(w http.ResponseWriter, r *http.Request) {
	server.InitServer()
	engine.ServeHTTP(w, r)
}
