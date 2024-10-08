package store

import (
	"bytes"
	"fmt"
	"grpcCource/pkg/models"
	"os"
	"sync"

	"github.com/google/uuid"
)

type ImageStore interface {
	Save(laptopID string, imageType string, imageData bytes.Buffer) (string, error)
}

type DickImageStore struct {
	mutex       sync.Mutex
	imageFolder string
	images      map[string]*models.ImageInfo
}

func NewDickImageStore(imageFolder string) *DickImageStore {
	store := &DickImageStore{
		imageFolder: imageFolder,
		images:      map[string]*models.ImageInfo{},
	}
	return store
}

func (store *DickImageStore) Save(laptopID string, imageType string, imageData bytes.Buffer) (string, error) {
	imageID := uuid.NewString()
	imagePath := fmt.Sprintf("%s/%s%s", store.imageFolder, imageID, imageType)
	file, err := os.Create(imagePath)
	if err != nil {
		return "", fmt.Errorf("Cannot create image file: %v", err)
	}

	_, err = imageData.WriteTo(file)
	if err != nil {
		return "", fmt.Errorf("Cannot wirte image to file: %v", err)
	}

	store.mutex.Lock()
	store.images[imagePath] = &models.ImageInfo{
		LaptopID: laptopID,
		Type:     imageType,
		Path:     imagePath,
	}
	store.mutex.Unlock()
	return imageID, nil
}
