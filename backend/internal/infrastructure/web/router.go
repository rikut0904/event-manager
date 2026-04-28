package web

import (
	"backend/internal/interface/handler"
	"net/http"
)

func NewRouter(healthHandler *handler.HealthHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", healthHandler.Handle)
	return mux
}
