// Package users
// This package conains the user service files
package users

import (
	"context"
	"log"
	"log/slog"

	"github.com/AzmainMahtab/go-chi-hex/internal/domain"
	"github.com/AzmainMahtab/go-chi-hex/internal/ports"
	"github.com/google/uuid"
)

type service struct {
	repo   ports.UserRepository
	hasher ports.PasswordHasher
}

func NewUserService(repo ports.UserRepository, hasher ports.PasswordHasher) ports.UserService {
	return &service{
		repo:   repo,
		hasher: hasher,
	}
}

// RegisterUser takes a domain.User and registers a user
func (s *service) RegisterUser(ctx context.Context, req domain.User) (*domain.User, error) {

	conflict, err := s.repo.CheckConflict(ctx, req.UserName, req.Email, req.Phone)
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
	hashedPass, err := s.hasher.Hash(req.Password)
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

	if err := s.repo.Create(ctx, &req); err != nil {
		log.Printf("Service: Create user error: %v", err)
		return nil, err
	}

	return &req, nil
}

func (s *service) ListUsers(ctx context.Context, filters domain.UserFilter) ([]*domain.User, error) {
	// showDeleted is false here because this is for "active" users
	users, err := s.repo.ReadAll(ctx, filters)
	if err != nil {
		log.Printf("Service: ReadAll error: %v", err)
		return nil, &domain.AppError{
			Code:    domain.CodeInternal,
			Message: "User List: Failed",
			Err:     err,
		}
	}
	return users, nil
}

func (s *service) GetUser(ctx context.Context, id string) (*domain.User, error) {
	u, err := s.repo.ReadOne(ctx, id)
	if err != nil {
		if u == nil {
			return nil, &domain.AppError{
				Code:    domain.CodeNotFound,
				Message: "User not found",
				Err:     err,
			}
		}
		return nil, &domain.AppError{
			Code:    domain.CodeInternal,
			Message: "Something went wrong in user fetching",
			Err:     err,
		}
	}
	return u, nil
}

func (s *service) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	u, err := s.repo.ReadByEmail(ctx, email)
	if err != nil {
		if u == nil {
			return nil, &domain.AppError{
				Code:    domain.CodeValidation,
				Message: "User not found",
				Err:     err,
			}
		}
		return nil, &domain.AppError{
			Code:    domain.CodeValidation,
			Message: "Something went wrong",
			Err:     err,
		}

	}
	return u, nil
}

func (s *service) UpdateUser(ctx context.Context, updates domain.UserUpdate) (*domain.User, error) {
	// Check if user exists first (Optional, but good for business logic)
	_, err := s.repo.ReadOne(ctx, updates.UUID)
	if err != nil {
		return nil, &domain.AppError{
			Code:    domain.CodeNotFound,
			Message: "Resource Not found",
			Err:     err,
		}
	}

	// Perform the partial update
	if err := s.repo.Update(ctx, updates); err != nil {
		slog.Error("Update err:", "err", err)
		return nil, &domain.AppError{
			Code:    domain.CodeInternal,
			Message: "Action could not be performed",
			Err:     err,
		}
	}

	// Return the fresh user data
	return s.repo.ReadOne(ctx, updates.UUID)
}

func (s *service) RemoveUser(ctx context.Context, id string) error {
	_, err := s.repo.ReadOne(ctx, id)
	if err != nil {
		return &domain.AppError{
			Code:    domain.CodeNotFound,
			Message: "Resource Not found",
			Err:     err,
		}
	}

	err = s.repo.SoftDelete(ctx, id)
	if err != nil {
		return &domain.AppError{
			Code:    domain.CodeInternal,
			Message: "Action could not be performed",
			Err:     err,
		}
	}
	return nil

}

func (s *service) RestoreUser(ctx context.Context, id string) (*domain.User, error) {
	_, err := s.repo.ReadOne(ctx, id)
	if err != nil {
		return nil, &domain.AppError{
			Code:    domain.CodeNotFound,
			Message: "Resource not found",
			Err:     err,
		}
	}

	err = s.repo.Restore(ctx, id)
	if err != nil {
		return nil, &domain.AppError{
			Code:    domain.CodeInternal,
			Message: "Action could not be performed",
			Err:     err,
		}
	}

	return s.repo.ReadOne(ctx, id)

}

func (s *service) GetTrashedUsers(ctx context.Context, filters domain.UserFilter) ([]*domain.User, error) {
	// We pass an empty filter map to get all trashed users for now
	usr, err := s.repo.Trash(ctx, filters)
	if err != nil {
		return nil, &domain.AppError{
			Code:    domain.CodeInternal,
			Message: "Action could not be perforemed",
			Err:     err,
		}
	}

	return usr, nil
}

func (s *service) PermanentlyDeleteUser(ctx context.Context, id string) error {
	usr, err := s.repo.ReadOne(ctx, id)
	if err != nil {
		return &domain.AppError{
			Code:    domain.CodeNotFound,
			Message: "Resource not found",
			Err:     err,
		}
	}

	if usr.UserStatus != "inactive" {
		return &domain.AppError{
			Code:    domain.CodeValidation,
			Message: "Can not delete active user",
			Err:     err,
		}
	}

	err = s.repo.Prune(ctx, id)
	if err != nil {
		return &domain.AppError{
			Code:    domain.CodeInternal,
			Message: "Action could not be performed",
			Err:     err,
		}
	}

	return nil
}
