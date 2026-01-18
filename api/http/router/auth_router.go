// Package routes
// this conains the auth routes
package routes

import (
	"net/http"

	"github.com/AzmainMahtab/go-chi-hex/api/http/handlers"
	"github.com/go-chi/chi/v5"
)

func authRouter(ah *handlers.AuthHandler) http.Handler {
	r := chi.NewRouter()

	r.Post("/login", ah.Login) // Get the token pair

	return r
}
