package store

import (
	"context"
	"errors"
	"fmt"
	"grpcCource/pkg/pb"
	"sync"

	"github.com/jinzhu/copier"
)

var ErrAlreadyExists = errors.New("record already exists")

type LaptopStore interface {
	Save(laptop *pb.Laptop) error
	Find(id string) (*pb.Laptop, error)
	Search(ctx context.Context, filter *pb.Filter, found func(laptop *pb.Laptop) error) error
}

type InMemoryLaptopStore struct {
	data  map[string]*pb.Laptop
	mutex sync.Mutex
}

func NewInMemoryLaptopStore() *InMemoryLaptopStore {
	return &InMemoryLaptopStore{
		data: make(map[string]*pb.Laptop),
	}
}

func (store *InMemoryLaptopStore) Save(laptop *pb.Laptop) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	if store.data[laptop.Id] != nil {
		return ErrAlreadyExists
	}
	other := &pb.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil {
		return fmt.Errorf("Cannot copy laptop: %v", err)
	}
	store.data[laptop.Id] = laptop
	return nil
}

func (store *InMemoryLaptopStore) Find(id string) (*pb.Laptop, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	laptop, ok := store.data[id]
	if !ok {
		return nil, nil
	}
	other, err := deepCopy[pb.Laptop](laptop)
	if err != nil {
		return nil, fmt.Errorf("Cannot copy laptop: %v", err)
	}
	return other, nil
}

// type DBLaptopStore struct{}
func (store *InMemoryLaptopStore) Search(
	ctx context.Context,
	filter *pb.Filter,
	found func(laptop *pb.Laptop) error,
) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	for _, laptop := range store.data {
		if ctx.Err() == context.DeadlineExceeded || ctx.Err() == context.Canceled {
			return ctx.Err()
		}
		if isQualified(filter, laptop) {
			other, err := deepCopy[pb.Laptop](laptop)
			if err != nil {
				return err
			}
			err = found(other)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func isQualified(filter *pb.Filter, laptop *pb.Laptop) bool {
	if laptop.GetPriceUsd() > filter.GetMaxPriceUsd() {
		return false
	}
	if toBit(filter.GetMinRam()) > toBit(laptop.GetRam()) {
		return false
	}
	if laptop.GetCpu().GetNumberCores() < filter.GetMinCpuCores() {
		return false
	}
	if laptop.GetCpu().GetMinGhz() < filter.MinCpuGhz {
		return false
	}
	return true
}

func deepCopy[T any](i *T) (*T, error) {
	other := new(T)
	err := copier.Copy(other, i)
	if err != nil {
		return nil, err
	}
	return other, nil
}

func toBit(m *pb.Memory) uint64 {
	switch m.Unit {
	case pb.Memory_BIT:
		return m.Value
	case pb.Memory_BYTE:
		return m.Value << 3
	case pb.Memory_KILOBYTE:
		return m.Value << 13
	case pb.Memory_MEGABYTE:
		return m.Value << 23
	case pb.Memory_GIGABYTE:
		return m.Value << 33
	case pb.Memory_TETRABYTE:
		return m.Value << 43
	default:
		return 0
	}
}
