package service

import (
	"errors"
	"grpcCource/pkg/models"
	"sync"
)

var ErrNotFound = errors.New("Not found in system")

type InMemoryUserStore struct {
	mutext sync.RWMutex
	users  map[string]*models.User
}

func NewInMemoryUserStore() *InMemoryUserStore {
	store := &InMemoryUserStore{users: map[string]*models.User{}}
	return store
}

func (store *InMemoryUserStore) Add(user *models.User) error {
	store.mutext.Lock()
	defer store.mutext.Unlock()
	_, ok := store.users[user.Username]
	if ok {
		return ErrAlreadyExists
	}
	store.users[user.Username] = user
	return nil
}

func (store *InMemoryUserStore) Find(userName string) (*models.User, error) {
	u, ok := store.users[userName]
	if !ok {
		return nil, nil
	}
	return u, nil
}
