// Package auth
// This package handles auth service
package auth

import (
	"context"
	"log"
	"time"

	"github.com/AzmainMahtab/go-chi-hex/internal/domain"
	"github.com/AzmainMahtab/go-chi-hex/internal/ports"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	repo          ports.UserRepository
	tokenProvider ports.TokenProvider
	cache         ports.CacheRepo
	hasher        ports.PasswordHasher
	auditPub      ports.AuditPublisher
}

func NewAuthService(ur ports.UserRepository, tp ports.TokenProvider, c ports.CacheRepo, h ports.PasswordHasher, ap ports.AuditPublisher) ports.AuthService {
	return &authService{
		repo:          ur,
		tokenProvider: tp,
		cache:         c,
		hasher:        h,
		auditPub:      ap,
	}
}

func (a *authService) Register(ctx context.Context, req domain.User) (*domain.User, error) {

	conflict, err := a.repo.CheckConflict(ctx, req.UserName, req.Email, req.Phone)
	if err != nil {
		return nil, &domain.AppError{
			Code:    domain.CodeInternal,
			Message: "Database check failed",
			Err:     err,
		}
	}

	if len(conflict) > 0 {
		return nil, &domain.AppError{
			Code:    domain.CodeConflict,
			Message: "User register failed: Conflicting values",
			Errors:  conflict,
		}
	}

	// Hashing the Password
	hashedPass, err := a.hasher.Hash(req.Password)
	if err != nil {
		return nil, &domain.AppError{
			Code:    domain.CodeInternal,
			Message: "Password could not be hashedPass",
			Err:     err,
		}
	}

	//Seting hashed password and generating UUID V7
	req.Password = hashedPass
	req.UserRole = "user"

	newUUID, _ := uuid.NewV7()
	req.UUID = newUUID.String()

	if err := a.repo.Create(ctx, &req); err != nil {
		log.Printf("Service: Create user error: %v", err)
		return nil, err
	}

	return &req, nil
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

	eventUUID, _ := uuid.NewV7()

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(login.Password)); err != nil {

		a.auditPub.Publish(ctx, domain.Audit{
			UUID:      eventUUID.String(),
			EventType: "USER_LOGIN",
			ActorID:   u.UUID,
			Payload: map[string]any{
				"email":  u.Email,
				"status": "Failed",
			},
		})

		return domain.Tokenpair{}, &domain.AppError{
			Code:    domain.CodeValidation,
			Message: "One or more wrong credential",
			Err:     err,
		}
	}

	a.auditPub.Publish(ctx, domain.Audit{
		UUID:      eventUUID.String(),
		EventType: "USER_LOGIN",
		ActorID:   u.UUID,
		Payload: map[string]any{
			"email":  u.Email,
			"status": "Success",
		},
	})

	return a.tokenProvider.GenerateTokenPair(u)
}

func (a *authService) Logout(ctx context.Context, refreshToken string, accessClaims domain.UserClaims) error {

	exists, err := a.cache.Exists(ctx, "blacklist:refresh:"+refreshToken)
	if err != nil {
		return &domain.AppError{
			Code:    domain.CodeInternal,
			Message: "Something happened",
			Err:     err,
		}
	}
	if exists {
		return &domain.AppError{
			Code:    domain.CodeUauthorized,
			Message: "BAAAAD TOOOKKKEEEENN",
			Err:     err,
		}

	}

	//  Verify the refresh token to get ITS specific expiration
	refreshClaims, err := a.tokenProvider.VerifyToken(refreshToken)
	if err != nil {
		return &domain.AppError{
			Code:    domain.CodeUauthorized,
			Message: "BAAAAAD TOKEN",
			Err:     err,
		}
	}

	// Make sure the person owning the Access Token is the one owning the Refresh Token
	if accessClaims.UserID != refreshClaims.UserID {
		return &domain.AppError{
			Code:    domain.CodeUauthorized,
			Message: "Bad token",
			Err:     err,
		}
	}

	//  Get expiration time from the REFRESH claims
	expTime := time.Unix(refreshClaims.Expires, 0)
	ttl := time.Until(expTime)

	if ttl <= 0 {
		return nil
	}

	//  Blacklist the token string
	return a.cache.Set(ctx, "blacklist:refresh:"+refreshToken, "revoked", ttl)
}
func (a *authService) Rotate(ctx context.Context, refreshToken string) (domain.Tokenpair, error) {
	// Checking if token already exist
	blackList, err := a.cache.Exists(ctx, "blacklist:refresh:"+refreshToken)
	if err != nil {
		return domain.Tokenpair{}, &domain.AppError{
			Code:    domain.CodeInternal,
			Message: "Something happened",
			Err:     err,
		}
	}

	if blackList {
		return domain.Tokenpair{}, &domain.AppError{
			Code:    domain.CodeUauthorized,
			Message: "Bad token",
			Err:     err,
		}
	}

	claims, err := a.tokenProvider.VerifyToken(refreshToken)
	if err != nil {
		return domain.Tokenpair{}, &domain.AppError{
			Code:    domain.CodeInternal,
			Message: "Something happened",
			Err:     err,
		}
	}

	usr, err := a.repo.ReadOne(ctx, claims.UserID)
	if err != nil {
		return domain.Tokenpair{}, &domain.AppError{
			Code:    domain.CodeInternal,
			Message: "Something happened",
			Err:     err,
		}
	}

	if usr.UserStatus != "active" {
		return domain.Tokenpair{}, &domain.AppError{
			Code:    domain.CodeInternal,
			Message: "Uauthorized !",
			Err:     err,
		}
	}

	newToken, err := a.tokenProvider.GenerateTokenPair(usr)
	if err != nil {
		return domain.Tokenpair{}, &domain.AppError{
			Code:    domain.CodeInternal,
			Message: "Something happened",
			Err:     err,
		}
	}

	expTime := time.Unix(claims.Expires, 0)
	ttl := time.Until(expTime)

	if ttl > 0 {
		if err := a.cache.Set(ctx, "blacklist:refresh:"+refreshToken, "Rotated", ttl); err != nil {
			return domain.Tokenpair{}, &domain.AppError{
				Code:    domain.CodeInternal,
				Message: "Something happened",
				Err:     err,
			}

		}
	}

	return newToken, nil

}
