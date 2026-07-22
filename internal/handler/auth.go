package handler

import (
	"errors"
	"net/http"

	"github.com/RahulRazj/go-crud-auth/internal/errs"
	"github.com/RahulRazj/go-crud-auth/internal/model"
	"github.com/RahulRazj/go-crud-auth/internal/repository"
	"github.com/RahulRazj/go-crud-auth/internal/server"
	"github.com/RahulRazj/go-crud-auth/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	Handler
	auth     *service.AuthService
	validate *validator.Validate
}

func NewAuthHandler(s *server.Server, users *repository.UserRepository) *AuthHandler {
	return &AuthHandler{
		Handler:  NewHandler(s),
		auth:     service.NewAuthService(users),
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

	user, err := h.auth.Signup(c.Request().Context(), input)
	if errors.Is(err, repository.ErrEmailExists) {
		return errs.NewBadRequestError("Email is already registered", false, nil, nil, nil)
	}
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, model.AuthResponse{User: user})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var input model.LoginRequest
	if err := c.Bind(&input); err != nil {
		return errs.NewBadRequestError("Invalid request body", false, nil, nil, nil)
	}
	if err := h.validate.Struct(input); err != nil {
		return errs.ValidationError(err)
	}

	user, err := h.auth.Login(c.Request().Context(), input)
	if errors.Is(err, service.ErrInvalidCredentials) {
		return errs.NewUnauthorizedError("Invalid email or password", false)
	}
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, model.AuthResponse{User: user})
}
