package service

import (
	"sync"
)

type RatingStore interface {
	Add(laptopID string, score float64) (*Rating, error)
}
type Rating struct {
	Count uint32
	Sum   float64
}

type InMemoryRatingStore struct {
	mutext  sync.RWMutex
	ratings map[string]*Rating
}

func NewInMemoryRatingStore() *InMemoryRatingStore {
	store := &InMemoryRatingStore{
		ratings: map[string]*Rating{},
	}
	return store
}

func (store *InMemoryRatingStore) Add(laptopID string, score float64) (*Rating, error) {
	store.mutext.Lock()
	defer store.mutext.Unlock()

	rating, ok := store.ratings[laptopID]
	if !ok {
		rating = &Rating{
			Count: 1,
			Sum:   score,
		}
		store.ratings[laptopID] = rating
		return rating, nil
	}

	rating.Sum += score
	rating.Count += 1
	store.ratings[laptopID] = rating
	copyRating, err := deepCopy[Rating](rating)
	if err != nil {
		return nil, err
	}
	return copyRating, nil
}
