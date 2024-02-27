package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"todoPchi"
)

type CommentStore struct {
	*mongo.Collection
}

func (s *CommentStore) Comments(id primitive.ObjectID) ([]todoPchi.Comment, error) {
	var result []todoPchi.Comment
	filter := bson.D{{"todo_id", id}}
	cursor, err := s.Collection.Find(context.TODO(), filter)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var comment todoPchi.Comment
		if err := cursor.Decode(&comment); err != nil {
			log.Println(err)
			return nil, err
		}
		result = append(result, comment)
	}
	if err := cursor.Err(); err != nil {
		log.Println(err)
		return nil, err
	}
	return result, nil
}

func (s *CommentStore) CommentCreate(id primitive.ObjectID, t *todoPchi.Comment) error {
	comment := bson.D{
		//{"todo_id", t.TodoID},
		{"description", t.Description},
		{"date", t.Date},
	}

	_, err := s.InsertOne(context.TODO(), comment)
	if err != nil {
		log.Printf("Error inserting comment: %v\n", err)
		return err
	}

	log.Printf("Inserted new comment with ID: %v\n", t.ID.Hex())
	return nil
}
