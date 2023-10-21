package service

import (
	"context"
	"grpcCource/pkg/models"
	"grpcCource/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoUserStore struct {
	collection *mongo.Collection
}

func NewMongoUserStore(db *mongo.Database) *MongoUserStore {
	userStore := &MongoUserStore{
		collection: db.Collection("users"),
	}
	return userStore
}

func (store *MongoUserStore) Add(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := store.collection.InsertOne(ctx, user)
	return err
}
func (store *MongoUserStore) Find(userName string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	filter := bson.M{
		"user_name": userName,
	}
	return utils.FindOne[models.User](ctx, store.collection, filter)
}
