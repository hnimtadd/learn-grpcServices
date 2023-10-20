package store

import (
	"grpcCource/pkg/models"
)

type RatingStore interface {
	Add(laptopID string, score float64) (*models.Rating, error)
}
