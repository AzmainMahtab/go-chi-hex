// Package auth
// This package handles auth service
package auth

import (
	"context"

	"github.com/AzmainMahtab/go-chi-hex/internal/domain"
	"github.com/AzmainMahtab/go-chi-hex/internal/ports"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	repo          ports.UserRepository
	tokenProvider ports.TokenProvider
}

func NewAuthService(ur ports.UserRepository, tp ports.TokenProvider) ports.AuthService {
	return &authService{
		repo:          ur,
		tokenProvider: tp,
	}
}

func (a *authService) Login(ctx context.Context, login domain.AuthLogin) (domain.Tokenpair, error) {
	u, err := a.repo.ReadByEmail(ctx, login.Email)
	if err != nil {
		return domain.Tokenpair{}, &domain.AppError{
			Code:    domain.CodeValidation,
			Message: "One or more wrong credential",
			Err:     err,
		}
	}

	if u.UserStatus != "active" {
		return domain.Tokenpair{}, &domain.AppError{
			Code:    domain.CodeValidation,
			Message: "Account suspended or inactive. contact admin",
			Err:     err,
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(login.Password)); err != nil {
		return domain.Tokenpair{}, &domain.AppError{
			Code:    domain.CodeValidation,
			Message: "One or more wrong credential",
			Err:     err,
		}
	}

	return a.tokenProvider.GenerateTokenPair(u)
}
