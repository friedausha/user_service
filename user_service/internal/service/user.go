package service

import (
	"context"
	"database/sql"
	"errors"
	"git.garena.com/frieda.hasanah/user_service/internal/data/cache"
	"git.garena.com/frieda.hasanah/user_service/internal/model"
	"git.garena.com/frieda.hasanah/user_service/utils/hash"
)

type service struct {
	repo  model.IUserRepository
	cache *cache.UserCache
}

// Register user with unique username
func (s service) Register(ctx context.Context, user model.User) (model.User, error) {
	// Retrieve from cache first
	cachedUser, found := s.cache.Get(user.Username)
	if found {
		return cachedUser, nil
	}

	u, err := s.repo.GetUserByUsername(ctx, user.Username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return model.User{}, errors.New("error fetching user db")
	}
	if u.Username != "" {
		return model.User{}, errors.New("user already exists")
	}

	// hash user's password
	hashedPassword, err := hash.EncryptPassword(user.Password)
	if err != nil {
		return model.User{}, errors.New("failed to hash password")
	}
	user.Password = hashedPassword

	id, err := s.repo.Create(ctx, user)
	if err != nil {
		return model.User{}, err
	}

	user.ID = id
	s.cache.Set(user)
	return user, nil
}

// NewService will initialize the implementations of auth service
func NewService(
	userRepo model.IUserRepository,
	cache *cache.UserCache,
) model.IUserService {
	return &service{
		repo:  userRepo,
		cache: cache,
	}
}
