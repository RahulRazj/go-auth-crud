package service

import (
	"testing"
	"time"

	"github.com/RahulRazj/go-crud-auth/internal/config"
	"github.com/RahulRazj/go-crud-auth/internal/model"
	"github.com/google/uuid"
)

func testJWTConfig() config.JWTConfig {
	return config.JWTConfig{
		Secret:               "test-secret-key-12345",
		AccessTokenTTLMin:   15,
		RefreshTokenTTLHours: 168,
	}
}

func TestGenerateAndValidateTokenPair(t *testing.T) {
	jwtSvc := NewJWTService(testJWTConfig())

	userID := uuid.New()
	user := model.User{
		ID:        userID,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
	}

	beforeGen := time.Now()
	tokenPair, err := jwtSvc.GenerateTokenPair(user)
	afterGen := time.Now()

	if err != nil {
		t.Fatalf("GenerateTokenPair failed: %v", err)
	}

	if tokenPair.AccessToken == "" {
		t.Errorf("expected non-empty access token")
	}

	if tokenPair.RefreshToken == "" {
		t.Errorf("expected non-empty refresh token")
	}

	// Verify expiry times
	expectedAccessExpiry := beforeGen.Add(15 * time.Minute)
	if tokenPair.AccessTokenExpiresAt.Before(expectedAccessExpiry) || tokenPair.AccessTokenExpiresAt.After(afterGen.Add(15*time.Minute)) {
		t.Errorf("unexpected AccessTokenExpiresAt: %v", tokenPair.AccessTokenExpiresAt)
	}

	expectedRefreshExpiry := beforeGen.Add(168 * time.Hour)
	if tokenPair.RefreshTokenExpiresAt.Before(expectedRefreshExpiry) || tokenPair.RefreshTokenExpiresAt.After(afterGen.Add(168*time.Hour)) {
		t.Errorf("unexpected RefreshTokenExpiresAt: %v", tokenPair.RefreshTokenExpiresAt)
	}

	// Validate Access Token
	accessClaims, err := jwtSvc.ValidateToken(tokenPair.AccessToken, TokenTypeAccess)
	if err != nil {
		t.Fatalf("ValidateToken for access token failed: %v", err)
	}

	if accessClaims.UserID != userID {
		t.Errorf("expected UserID %v, got %v", userID, accessClaims.UserID)
	}

	if accessClaims.Email != "john@example.com" {
		t.Errorf("expected Email john@example.com, got %v", accessClaims.Email)
	}

	if accessClaims.TokenType != TokenTypeAccess {
		t.Errorf("expected TokenTypeAccess, got %v", accessClaims.TokenType)
	}

	// Validate Refresh Token
	refreshClaims, err := jwtSvc.ValidateToken(tokenPair.RefreshToken, TokenTypeRefresh)
	if err != nil {
		t.Fatalf("ValidateToken for refresh token failed: %v", err)
	}

	if refreshClaims.UserID != userID {
		t.Errorf("expected UserID %v, got %v", userID, refreshClaims.UserID)
	}

	if refreshClaims.Email != "john@example.com" {
		t.Errorf("expected Email john@example.com, got %v", refreshClaims.Email)
	}

	if refreshClaims.TokenType != TokenTypeRefresh {
		t.Errorf("expected TokenTypeRefresh, got %v", refreshClaims.TokenType)
	}
}

func TestValidateTokenMismatch(t *testing.T) {
	jwtSvc := NewJWTService(testJWTConfig())

	user := model.User{
		ID:    uuid.New(),
		Email: "test@example.com",
	}

	tokenPair, err := jwtSvc.GenerateTokenPair(user)
	if err != nil {
		t.Fatalf("GenerateTokenPair failed: %v", err)
	}

	// Try validating Access Token as Refresh Token -> should fail
	_, err = jwtSvc.ValidateToken(tokenPair.AccessToken, TokenTypeRefresh)
	if err == nil {
		t.Errorf("expected error when validating access token as refresh token, got nil")
	}

	// Try validating Refresh Token as Access Token -> should fail
	_, err = jwtSvc.ValidateToken(tokenPair.RefreshToken, TokenTypeAccess)
	if err == nil {
		t.Errorf("expected error when validating refresh token as access token, got nil")
	}
}

func TestValidateInvalidToken(t *testing.T) {
	jwtSvc := NewJWTService(testJWTConfig())

	_, err := jwtSvc.ValidateToken("invalid.jwt.token", TokenTypeAccess)
	if err == nil {
		t.Errorf("expected error for invalid token, got nil")
	}
}
