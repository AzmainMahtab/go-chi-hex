// Package routes
// this conains the auth routes
package routes

import (
	"net/http"

	"github.com/AzmainMahtab/go-chi-hex/api/http/handlers"
	"github.com/AzmainMahtab/go-chi-hex/api/http/middleware"
	"github.com/AzmainMahtab/go-chi-hex/internal/ports"
	"github.com/go-chi/chi/v5"
)

func authRouter(ah *handlers.AuthHandler, tokenProvider ports.TokenProvider) http.Handler {
	r := chi.NewRouter()

	//  PUBLIC ROUTES No Middlewar
	r.Post("/register", ah.Register)
	r.Post("/login", ah.Login)
	r.Post("/rotate", ah.Rotate)

	//  PROTECTED ROUTES
	r.With(middleware.AuthMiddleware(tokenProvider)).Post("/logout", ah.Logout)

	return r
}
