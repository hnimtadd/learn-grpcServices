package service

import (
	"grpcCource/pkg/models"
	"grpcCource/utils"
	"sync"
)

type InMemoryRatingStore struct {
	mutext  sync.RWMutex
	ratings map[string]*models.Rating
}

func NewInMemoryRatingStore() *InMemoryRatingStore {
	store := &InMemoryRatingStore{
		ratings: map[string]*models.Rating{},
	}
	return store
}

func (store *InMemoryRatingStore) Add(laptopID string, score float64) (*models.Rating, error) {
	store.mutext.Lock()
	defer store.mutext.Unlock()

	rating, ok := store.ratings[laptopID]
	if !ok {
		rating = &models.Rating{
			Count: 1,
			Sum:   score,
		}
		store.ratings[laptopID] = rating
		return rating, nil
	}

	rating.Sum += score
	rating.Count += 1
	store.ratings[laptopID] = rating
	copyRating, err := utils.DeepCopy[models.Rating](rating)
	if err != nil {
		return nil, err
	}
	return copyRating, nil
}
