package model

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID    `json:"id" db:"id" validate:"required"`
	Email     string       `json:"email" db:"email" validate:"required,email"`
	Username  string       `json:"username" db:"username" validate:"required"`
	FullName  string       `json:"full_name" db:"full_name" validate:"required"`
	Password  string       `json:"password" validate:"required"`
	Token     *string      `json:"token"`
	CreatedAt sql.NullTime `json:"created_at" db:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at" db:"updated_at"`
}

type IUserRepository interface {
	Create(ctx context.Context, user User) (uuid.UUID, error)
	GetUserByUsername(ctx context.Context, username string) (User, error)
	UpdateUserToken(ctx context.Context, userID, token string) error
}

type IUserService interface {
	Register(ctx context.Context, user User) (User, error)
}

type IAuthService interface {
	VerifyToken(tokenString string) (bool, error)
	Login(ctx context.Context, username, password string) (User, int64, error)
}
