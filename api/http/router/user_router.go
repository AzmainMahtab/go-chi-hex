// Package routes
// this contains the user routes
package routes

import (
	"net/http"

	"github.com/AzmainMahtab/go-chi-hex/api/http/handlers"
	"github.com/go-chi/chi/v5"
)

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
