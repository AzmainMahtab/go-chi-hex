// Package routes
// routes.go combinds all the other routs at one place
package routes

import (
	"net/http"

	"github.com/AzmainMahtab/docpad/api/http/handlers"
	_ "github.com/AzmainMahtab/docpad/docs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
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
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", deps.HealthH.HealthCheck)
	})

	// --- Static Handler for /docs/* ---
	fileServer := http.FileServer(http.Dir("./docs"))
	r.Handle("/docs/*", http.StripPrefix("/docs", fileServer))

	// Swagger here
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/docs/swagger.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
	))

	return r
}
