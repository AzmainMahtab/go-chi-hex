// Package users
// This package conains the user service files
package users

import (
	"context"
	"errors"
	"log"

	"github.com/AzmainMahtab/docpad/api/http/dto"
	"github.com/AzmainMahtab/docpad/internal/domain"
	"github.com/AzmainMahtab/docpad/internal/ports"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) ports.UserService {
	return &service{repo: repo}
}

func (s *service) RegisterUser(ctx context.Context, req dto.RegisterUserRequest) (*dto.UserResponse, error) {
	pass, err := s.hashPassword(req.Password)
	if err != nil {
		log.Printf("hashing problem: %v", err)
		return nil, err
	}

	userDomain := &domain.User{
		UserName: req.UserName,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: pass,
	}

	if err := s.repo.Create(ctx, userDomain); err != nil {
		log.Printf("Can not create user: %v", err)
		return nil, err
	}

	return s.mapToResponse(userDomain), nil
}

func (s *service) ListUsers(ctx context.Context, filters map[string]any) ([]*dto.UserResponse, error) {
	return nil, nil
}

func (s *service) GetUser(ctx context.Context, id int) (*dto.UserResponse, error) {
	u, err := s.repo.ReadOne(ctx, id)
	if err != nil {
		log.Printf("Readone repo error: %v", err)
		return nil, err
	}
	if u == nil {
		return nil, errors.New("user not found")
	}

	return s.mapToResponse(u), nil
}

func (s *service) UpdateUser(ctx context.Context, id int, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	return nil, nil
}

func (s *service) RemoveUser(ctx context.Context, id int) error {
	return nil
}

func (s *service) GetTrashedUsers(ctx context.Context) ([]*dto.UserResponse, error) {
	return nil, nil
}

func (s *service) PermanentlyDeleteUser(ctx context.Context, id int) error {
	return nil
}

// --- PRIVATE HELPERS FOR MAPPING ---

func (s *service) mapToResponse(u *domain.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:         u.ID,
		UserName:   u.UserName,
		Email:      u.Email,
		Phone:      u.Phone,
		UserStatus: u.UserStatus,
		CreatedAt:  u.CreatedAt,
	}
}

// func (s *service) mapSliceToResponse(users []*domain.User) []*dto.UserResponse {
// 	res := make([]*dto.UserResponse, len(users))
// 	for i, u := range users {
// 		res[i] = s.mapToResponse(u)
// 	}
// 	return res
// }

//  --- PASSWORD HASHING HELPER ---//

func (s *service) hashPassword(pass string) (string, error) {
	byte, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return string(byte), err
}
