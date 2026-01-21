package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/AzmainMahtab/go-chi-hex/internal/ports"
	"github.com/AzmainMahtab/go-chi-hex/pkg/jsonutil"
)

// a custom type for context keys to avoid collisions
type contextKey string

const UserContextKey contextKey = "user_claims"

func AuthMiddleware(tokenProvider ports.TokenProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//  Get the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				jsonutil.WriteJSON(w, http.StatusUnauthorized, nil, nil, "Missing authorization header")
				return
			}

			//  Parse the Bearer token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				jsonutil.WriteJSON(w, http.StatusUnauthorized, nil, nil, "Invalid authorization format")
				return
			}

			tokenString := parts[1]

			//  Verify the token using your provider
			claims, err := tokenProvider.VerifyToken(tokenString)
			if err != nil {
				jsonutil.WriteJSON(w, http.StatusUnauthorized, nil, nil, "Invalid or expired token")
				return
			}

			//  Inject claims into the context and proceed
			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
