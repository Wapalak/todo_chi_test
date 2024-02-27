package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

func NewStore(url string) (*Store, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	// checking DB connection
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	err = client.Database("todo").CreateCollection(context.TODO(), "todo", nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	collection := client.Database("todo").Collection("todo")
	// Пока оставил так, а по идее надо поменять
	collection2 := client.Database("todo").Collection("todo")
	return &Store{
		TodoStore:    &TodoStore{Collection: collection},
		CommentStore: &CommentStore{Collection: collection2},
	}, err
}

type Store struct {
	*TodoStore
	*CommentStore
}
