package handler

import (
	"errors"
	"net/http"

	"github.com/RahulRazj/go-crud-auth/internal/errs"
	"github.com/RahulRazj/go-crud-auth/internal/middleware"
	"github.com/RahulRazj/go-crud-auth/internal/model"
	"github.com/RahulRazj/go-crud-auth/internal/repository"
	"github.com/RahulRazj/go-crud-auth/internal/server"
	"github.com/RahulRazj/go-crud-auth/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	Handler
	auth     *service.AuthService
	validate *validator.Validate
}

func NewAuthHandler(s *server.Server, users *repository.UserRepository, jwt *service.JWTService) *AuthHandler {
	return &AuthHandler{
		Handler:  NewHandler(s),
		auth:     service.NewAuthService(users, jwt),
		validate: validator.New(),
	}
}

func (h *AuthHandler) Signup(c echo.Context) error {
	var input model.SignupRequest
	if err := c.Bind(&input); err != nil {
		return errs.NewBadRequestError("Invalid request body", false, nil, nil, nil)
	}
	if err := h.validate.Struct(input); err != nil {
		return errs.ValidationError(err)
	}

	res, err := h.auth.Signup(c.Request().Context(), input)
	if errors.Is(err, repository.ErrEmailExists) {
		return errs.NewBadRequestError("Email is already registered", false, nil, nil, nil)
	}
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, res)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var input model.LoginRequest
	if err := c.Bind(&input); err != nil {
		return errs.NewBadRequestError("Invalid request body", false, nil, nil, nil)
	}
	if err := h.validate.Struct(input); err != nil {
		return errs.ValidationError(err)
	}

	res, err := h.auth.Login(c.Request().Context(), input)
	if errors.Is(err, service.ErrInvalidCredentials) {
		return errs.NewUnauthorizedError("Invalid email or password", false)
	}
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var input model.RefreshTokenRequest
	if err := c.Bind(&input); err != nil {
		return errs.NewBadRequestError("Invalid request body", false, nil, nil, nil)
	}
	if err := h.validate.Struct(input); err != nil {
		return errs.ValidationError(err)
	}

	res, err := h.auth.RefreshToken(c.Request().Context(), input.RefreshToken)
	if errors.Is(err, service.ErrInvalidToken) {
		return errs.NewUnauthorizedError("Invalid or expired refresh token", false)
	}
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) Me(c echo.Context) error {
	userID, ok := c.Get(middleware.ContextUserIDKey).(uuid.UUID)
	if !ok {
		return errs.NewUnauthorizedError("Unauthorized", false)
	}

	user, err := h.auth.GetMe(c.Request().Context(), userID)
	if errors.Is(err, service.ErrInvalidToken) {
		return errs.NewUnauthorizedError("User not found", false)
	}
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, model.AuthResponse{User: user})
}

