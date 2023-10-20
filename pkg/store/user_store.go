package store

import "grpcCource/pkg/models"

type UserStore interface {
	Add(user *models.User) error
	Find(userName string) (*models.User, error)
}
