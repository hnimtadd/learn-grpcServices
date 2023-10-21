package service

import (
	"context"
	"grpcCource/pkg/models"
	"grpcCource/pkg/store"
	"grpcCource/utils"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoScoreStore struct {
	collection *mongo.Collection
}

func NewMongoScoreStore(db *mongo.Database) store.RatingStore {
	scoreStore := &MongoScoreStore{
		collection: db.Collection("scores"),
	}
	return scoreStore
}

type MongoRating struct {
	LaptopID string `bson:"laptop_id,omiempty"`
	Rating   struct {
		Count uint32  `bson:"count,omiempty"`
		Sum   float64 `bson:"sum,omiempty"`
	} `bson:"rating,omiempty"`
}

func (store *MongoScoreStore) Add(laptopID string, score float64) (*models.Rating, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	filter := bson.M{"laptop_id": laptopID}
	update := bson.M{"$inc": bson.M{"rating.count": 1, "rating.sum": score}}
	res := store.collection.FindOneAndUpdate(ctx, filter, update)
	if err := res.Err(); err != nil {
		rating := MongoRating{
			LaptopID: laptopID,
			Rating: struct {
				Count uint32  "bson:\"count,omiempty\""
				Sum   float64 "bson:\"sum,omiempty\""
			}{
				Count: 1,
				Sum:   score,
			},
		}
		_, err := store.collection.InsertOne(ctx, rating)
		if err != nil {
			return nil, err
		}
	}

	result, err := utils.FindOne[MongoRating](ctx, store.collection, filter)
	if err != nil {
		log.Println("cannot find object", err)
		return nil, err
	}

	log.Println("Result: ", result)
	rating := &models.Rating{
		Count: result.Rating.Count,
		Sum:   float64(result.Rating.Sum),
	}
	log.Println("Rating:", rating)
	return rating, nil
}
