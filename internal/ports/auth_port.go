// Package ports
// This one has the auth ports
package ports

import (
	"context"

	"github.com/AzmainMahtab/go-chi-hex/internal/domain"
)

type TokenProvider interface {
	GenerateTokenPair(User *domain.User) (domain.Tokenpair, error)
	VerifyToken(token string) (domain.UserClaims, error)
}

type AuthService interface {
	Login(ctx context.Context, login domain.AuthLogin) (domain.Tokenpair, error)
	Logout(ctx context.Context, refreshToken string, claims domain.UserClaims) error
	Rotate(ctx context.Context, refreshToken string) (domain.Tokenpair, error)
}
