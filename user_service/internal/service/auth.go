package service

import (
	"context"
	"git.garena.com/frieda.hasanah/user_service/internal/data/cache"
	"git.garena.com/frieda.hasanah/user_service/internal/model"
	"git.garena.com/frieda.hasanah/user_service/utils/hash"
	"git.garena.com/frieda.hasanah/user_service/utils/token"
	"github.com/pkg/errors"
)

type authService struct {
	repo  model.IUserRepository
	cache *cache.UserCache
}

func (s authService) Login(ctx context.Context, username, password string) (model.User, int64, error) {
	var u model.User
	var err error

	// Retrieve from cache first
	u, found := s.cache.Get(username)
	if !found {
		u, err = s.repo.GetUserByUsername(ctx, username)
		if err != nil {
			return model.User{}, 0, errors.New("error fetching user db")
		}
	}

	isPasswordValid := hash.CheckPasswordHash(u.Password, password)
	if !isPasswordValid {
		return model.User{}, 0, errors.New("incorrect password")
	}

	tkn, expIn, err := token.Generate(u.ID.String())
	err = s.repo.UpdateUserToken(ctx, u.ID.String(), tkn)
	if err != nil {
		return model.User{}, 0, errors.New("error updating user token")
	}
	u.Token = &tkn
	return u, expIn, nil
}

func (s authService) VerifyToken(tokenString string) (bool, error) {
	valid, err := token.VerifyAndCheckExpiration(tokenString)
	if err != nil || !valid {
		return false, err
	}
	return true, nil
}

// NewAuthService will initialize the implementations of auth service
func NewAuthService(
	userRepo model.IUserRepository,
	cache *cache.UserCache,
) model.IAuthService {
	return &authService{
		repo:  userRepo,
		cache: cache,
	}
}
