package service

import (
	"errors"
	"sync"
)

type UserStore interface {
	Add(user *User) error
	Find(userName string) (*User, error)
}

var ErrNotFound = errors.New("Not found in system")

type InMemoryUserStore struct {
	mutext sync.RWMutex
	users  map[string]*User
}

func NewInMemoryUserStore() *InMemoryUserStore {
	store := &InMemoryUserStore{users: map[string]*User{}}
	return store
}

func (store *InMemoryUserStore) Add(user *User) error {
	store.mutext.Lock()
	defer store.mutext.Unlock()
	_, ok := store.users[user.Username]
	if ok {
		return ErrAlreadyExists
	}
	store.users[user.Username] = user
	return nil
}

func (store *InMemoryUserStore) Find(userName string) (*User, error) {
	u, ok := store.users[userName]
	if !ok {
		return nil, nil
	}
	return u, nil
}
