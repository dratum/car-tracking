package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"auto-tracking/internal/api/handler"
	apimw "auto-tracking/internal/api/middleware"
)

func NewRouter(deviceHandler *handler.DeviceHandler, apiKey string) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Device routes (API-key auth)
	r.Route("/api/device", func(r chi.Router) {
		r.Use(apimw.APIKeyAuth(apiKey))
		r.Post("/location", deviceHandler.PostLocation)
		r.Post("/trip/start", deviceHandler.PostTripStart)
		r.Post("/trip/end", deviceHandler.PostTripEnd)
	})

	// Web API routes (JWT auth)
	r.Route("/api/v1", func(r chi.Router) {
		// TODO: POST /auth/login (public)
		// TODO: Protected group with JWT middleware
		// TODO: GET /trips
		// TODO: GET /trips/{id}
		// TODO: GET /trips/{id}/points
		// TODO: GET /stats
	})

	return r
}
