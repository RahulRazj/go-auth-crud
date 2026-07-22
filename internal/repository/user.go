package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/RahulRazj/go-crud-auth/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrUserNotFound = errors.New("user not found")
var ErrEmailExists = errors.New("email already exists")

type UserRepository struct {
	BaseRepository
}

func NewUserRepository(base BaseRepository) *UserRepository {
	return &UserRepository{BaseRepository: base}
}

func (r *UserRepository) Create(ctx context.Context, firstName, lastName, email, passwordHash string) (model.User, error) {
	user := model.User{}
	err := r.db.Pool.QueryRow(ctx, `
		INSERT INTO users (first_name, last_name, email, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id, first_name, last_name, email, password_hash, created_at, updated_at`, firstName, lastName, email, passwordHash).
		Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return model.User{}, ErrEmailExists
		}
		return model.User{}, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (model.User, error) {
	user := model.User{}
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, first_name, last_name, email, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1`, email).
		Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.User{}, ErrUserNotFound
	}
	if err != nil {
		return model.User{}, fmt.Errorf("find user by email: %w", err)
	}

	return user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (model.User, error) {
	user := model.User{}
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, first_name, last_name, email, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1`, id).
		Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.User{}, ErrUserNotFound
	}
	if err != nil {
		return model.User{}, fmt.Errorf("find user by id: %w", err)
	}

	return user, nil
}
