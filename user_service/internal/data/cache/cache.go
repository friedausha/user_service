package cache

import (
	"git.garena.com/frieda.hasanah/user_service/internal/model"
	"sync"
)

type UserCache struct {
	Cache map[string]model.User
	mux   sync.RWMutex
}

// NewUserCache creates a new UserCache
func NewUserCache() *UserCache {
	return &UserCache{
		Cache: make(map[string]model.User),
	}
}

// Set sets a user in the cache
func (uc *UserCache) Set(user model.User) {
	uc.mux.Lock()
	defer uc.mux.Unlock()
	uc.Cache[user.Username] = user
}

// Get retrieves a user from the cache based on ID
func (uc *UserCache) Get(username string) (model.User, bool) {
	uc.mux.RLock()
	defer uc.mux.RUnlock()
	user, ok := uc.Cache[username]
	return user, ok
}

// Delete removes a user from the cache based on ID
func (uc *UserCache) Delete(username string) {
	uc.mux.Lock()
	defer uc.mux.Unlock()
	delete(uc.Cache, username)
}
