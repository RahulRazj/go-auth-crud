package service

import (
	"context"
	"errors"
	"strings"

	"github.com/RahulRazj/go-crud-auth/internal/model"
	"github.com/RahulRazj/go-crud-auth/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthService struct {
	users *repository.UserRepository
	jwt   *JWTService
}

func NewAuthService(users *repository.UserRepository, jwt *JWTService) *AuthService {
	return &AuthService{
		users: users,
		jwt:   jwt,
	}
}

func (s *AuthService) Signup(ctx context.Context, input model.SignupRequest) (model.AuthResponse, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.AuthResponse{}, err
	}

	user, err := s.users.Create(ctx, input.FirstName, input.LastName, email, string(passwordHash))
	if err != nil {
		return model.AuthResponse{}, err
	}

	tokens, err := s.jwt.GenerateTokenPair(user)
	if err != nil {
		return model.AuthResponse{}, err
	}

	return model.AuthResponse{
		User:                  user,
		AccessToken:           tokens.AccessToken,
		RefreshToken:          tokens.RefreshToken,
		AccessTokenExpiresAt:  tokens.AccessTokenExpiresAt,
		RefreshTokenExpiresAt: tokens.RefreshTokenExpiresAt,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, input model.LoginRequest) (model.AuthResponse, error) {
	user, err := s.users.FindByEmail(ctx, strings.ToLower(strings.TrimSpace(input.Email)))
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return model.AuthResponse{}, ErrInvalidCredentials
		}
		return model.AuthResponse{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return model.AuthResponse{}, ErrInvalidCredentials
	}

	tokens, err := s.jwt.GenerateTokenPair(user)
	if err != nil {
		return model.AuthResponse{}, err
	}

	return model.AuthResponse{
		User:                  user,
		AccessToken:           tokens.AccessToken,
		RefreshToken:          tokens.RefreshToken,
		AccessTokenExpiresAt:  tokens.AccessTokenExpiresAt,
		RefreshTokenExpiresAt: tokens.RefreshTokenExpiresAt,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (model.AuthResponse, error) {
	claims, err := s.jwt.ValidateToken(refreshToken, TokenTypeRefresh)
	if err != nil {
		return model.AuthResponse{}, ErrInvalidToken
	}

	user, err := s.users.FindByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return model.AuthResponse{}, ErrInvalidToken
		}
		return model.AuthResponse{}, err
	}

	tokens, err := s.jwt.GenerateTokenPair(user)
	if err != nil {
		return model.AuthResponse{}, err
	}

	return model.AuthResponse{
		User:                  user,
		AccessToken:           tokens.AccessToken,
		RefreshToken:          tokens.RefreshToken,
		AccessTokenExpiresAt:  tokens.AccessTokenExpiresAt,
		RefreshTokenExpiresAt: tokens.RefreshTokenExpiresAt,
	}, nil
}

func (s *AuthService) GetMe(ctx context.Context, userID uuid.UUID) (model.User, error) {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return model.User{}, ErrInvalidToken
		}
		return model.User{}, err
	}

	return user, nil
}
