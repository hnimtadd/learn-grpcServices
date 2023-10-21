package utils

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewMongoDB(dsn string, database string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	opt := options.Client().ApplyURI(dsn)
	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, &readpref.ReadPref{}); err != nil {
		return nil, err
	}
	return client.Database(database), nil
}

func FindOne[T any](ctx context.Context, collection *mongo.Collection, filter any, opts ...*options.FindOneOptions) (*T, error) {
	log.Println(filter)
	cur := collection.FindOne(ctx, filter, opts...)
	if err := cur.Err(); err != nil {
		return nil, err
	}
	ele := new(T)
	if err := cur.Decode(ele); err != nil {
		return nil, err
	}
	log.Println(ele)
	return ele, nil
}

func FindMany[T any](ctx context.Context, collection *mongo.Collection, filter any, opt ...*options.FindOptions) ([]*T, error) {
	cursor, err := collection.Find(ctx, filter, opt...)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	res := []*T{}
	for cursor.Next(ctx) {
		ele := new(T)
		if err := cursor.Decode(ele); err != nil {
			return res, err
		}
		res = append(res, ele)
	}
	if err := cursor.Err(); err != nil {
		return res, err
	}
	return res, nil
}
