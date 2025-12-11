package routes

import (
	"net/http"

	"github.com/AzmainMahtab/docpad/api/http/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type RouterDependencies struct {
	HealthH *handlers.HealthHandler
}

func NewRouter(deps RouterDependencies) http.Handler {
	r := chi.NewRouter()

	// chi middleware stack
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	// Main router group
	r.Route("api/v1", func(r chi.Router) {
		r.Get("/health", deps.HealthH.HealthCheck)
	})

	return r
}
