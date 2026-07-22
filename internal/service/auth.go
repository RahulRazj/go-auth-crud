package service

import (
	"context"
	"errors"
	"strings"

	"github.com/RahulRazj/go-crud-auth/internal/model"
	"github.com/RahulRazj/go-crud-auth/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthService struct {
	users *repository.UserRepository
}

func NewAuthService(users *repository.UserRepository) *AuthService {
	return &AuthService{users: users}
}

func (s *AuthService) Signup(ctx context.Context, input model.SignupRequest) (model.User, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, err
	}

	return s.users.Create(ctx, input.FirstName, input.LastName, email, string(passwordHash))
}

func (s *AuthService) Login(ctx context.Context, input model.LoginRequest) (model.User, error) {
	user, err := s.users.FindByEmail(ctx, strings.ToLower(strings.TrimSpace(input.Email)))
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return model.User{}, ErrInvalidCredentials
		}
		return model.User{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return model.User{}, ErrInvalidCredentials
	}

	return user, nil
}
