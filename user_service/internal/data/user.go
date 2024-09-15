package data

import (
	"context"
	"git.garena.com/frieda.hasanah/user_service/internal/model"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// UserRepository defines the interface for user database operation.
type UserRepository struct {
	DB *sqlx.DB
}

// NewUserRepository will initialize the MySQL repository for users.
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

// Create inserts a new user into the database.
func (r *UserRepository) Create(ctx context.Context, user model.User) (uuid.UUID, error) {
	newUUID := uuid.New()
	query := sq.
		Insert("users").
		Columns("id", "full_name", "username", "email", "password").
		Values(newUUID.String(), user.FullName, user.Username, user.Email, user.Password).
		PlaceholderFormat(sq.Question)

	_, err := query.RunWith(r.DB).ExecContext(ctx)
	if err != nil {
		return uuid.Nil, err
	}

	return newUUID, nil
}

// UpdateUserToken updates the token for user when successfully login.
func (r *UserRepository) UpdateUserToken(ctx context.Context, userID, token string) error {
	query := sq.
		Update("users").
		Set("token", token).
		Where(sq.Eq{"id": userID}).
		PlaceholderFormat(sq.Question)

	_, err := query.RunWith(r.DB).ExecContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

// GetUserByUsername get user by username
func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	query := sq.Select("*").From("users").Where(sq.Eq{"username": username}).Limit(1)
	stmt, args, err := query.ToSql()
	if err != nil {
		return model.User{}, err
	}
	var user model.User
	err = r.DB.GetContext(ctx, &user, stmt, args...)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}
