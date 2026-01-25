// Package routes
// routes.go combinds all the other routs at one place
package routes

import (
	"net/http"

	"github.com/AzmainMahtab/go-chi-hex/api/http/handlers"
	"github.com/AzmainMahtab/go-chi-hex/api/http/middleware"

	_ "github.com/AzmainMahtab/go-chi-hex/docs"
	"github.com/AzmainMahtab/go-chi-hex/internal/ports"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type RouterDependencies struct {
	HealthH *handlers.HealthHandler
	UserH   *handlers.UserHandler
	AuthH   *handlers.AuthHandler
}

func NewRouter(deps RouterDependencies, tokenProvider ports.TokenProvider) http.Handler {
	r := chi.NewRouter()

	// chi middleware stack
	// r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.RequestID)
	r.Use(middleware.StructuredLogger)
	r.Use(chiMiddleware.Recoverer)

	// Main router group
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", deps.HealthH.HealthCheck)
		r.Mount("/user", userRouter(deps.UserH, tokenProvider))
		r.Mount("/auth", authRouter(deps.AuthH, tokenProvider))
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
