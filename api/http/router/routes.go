// Package routes
// routes.go combinds all the other routs at one place
package routes

import (
	"net/http"

	"github.com/AzmainMahtab/go-chi-hex/api/http/handlers"
	_ "github.com/AzmainMahtab/go-chi-hex/docs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type RouterDependencies struct {
	HealthH *handlers.HealthHandler
	UserH   *handlers.UserHandler
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
		r.Mount("/user", userRouter(deps.UserH))
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

func userRouter(uh *handlers.UserHandler) http.Handler {
	r := chi.NewRouter()

	// General User Routes
	r.Post("/", uh.Register) // POST /user
	r.Get("/", uh.List)      // GET /user

	// Special route for trashed users
	r.Get("/trash", uh.GetTrashed) // GET /user/trash

	// Specific User ID Routes
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", uh.GetByID)          // GET /user/{id}
		r.Patch("/", uh.Update)         // PATCH /user/{id}
		r.Delete("/", uh.Remove)        // DELETE /user/{id} (Soft Delete)
		r.Patch("/restore", uh.Restore) // PATCH /user/{id}/restore (restore user)
		r.Delete("/prune", uh.Prune)    // DELETE /user/{id}/prune (Permanent)
	})

	return r
}
