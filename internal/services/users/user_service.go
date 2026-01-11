// Package users
// This package conains the user service files
package users

import (
	"context"
	"errors"
	"log"

	"github.com/AzmainMahtab/go-chi-hex/internal/domain"
	"github.com/AzmainMahtab/go-chi-hex/internal/ports"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) ports.UserService {
	return &service{repo: repo}
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

	hashedPass, err := s.hashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	req.Password = hashedPass

	if err := s.repo.Create(ctx, &req); err != nil {
		log.Printf("Service: Create user error: %v", err)
		return nil, err
	}

	return &req, nil
}

func (s *service) ListUsers(ctx context.Context, filters map[string]any) ([]*domain.User, error) {
	// showDeleted is false here because this is for "active" users
	users, err := s.repo.ReadAll(ctx, filters, false)
	if err != nil {
		log.Printf("Service: ReadAll error: %v", err)
		return nil, err
	}
	return users, nil
}

func (s *service) GetUser(ctx context.Context, id int) (*domain.User, error) {
	u, err := s.repo.ReadOne(ctx, id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (s *service) UpdateUser(ctx context.Context, id int, updates map[string]any) (*domain.User, error) {
	// Check if user exists first (Optional, but good for business logic)
	_, err := s.repo.ReadOne(ctx, id)
	if err != nil {
		return nil, err
	}

	// Perform the partial update
	if err := s.repo.Update(ctx, id, updates); err != nil {
		return nil, err
	}

	// Return the fresh user data
	return s.repo.ReadOne(ctx, id)
}

func (s *service) RemoveUser(ctx context.Context, id int) error {
	return s.repo.SoftDelete(ctx, id)
}

func (s *service) RestoreUser(ctx context.Context, id int) (*domain.User, error) {
	if err := s.repo.Restore(ctx, id); err != nil {
		return nil, err
	}

	return s.repo.ReadOne(ctx, id)
}

func (s *service) GetTrashedUsers(ctx context.Context) ([]*domain.User, error) {
	// We pass an empty filter map to get all trashed users for now
	return s.repo.Trash(ctx, make(map[string]any))
}

func (s *service) PermanentlyDeleteUser(ctx context.Context, id int) error {
	return s.repo.Prune(ctx, id)
}

// --- PRIVATE HELPERS ---

func (s *service) hashPassword(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return string(bytes), err
}
