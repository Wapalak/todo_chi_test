package database

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
	"todoPchi"
)

type TodoStore struct {
	*mongo.Collection
}

func (s *TodoStore) Todos() ([]todoPchi.Todo, error) {
	var results []todoPchi.Todo
	cursor, err := s.Collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.Background()) {
		var todo todoPchi.Todo
		if err := cursor.Decode(&todo); err != nil {
			log.Fatal(err)
		}

		// Directly parse the date using the expected format

		todo.Date = todo.Date.Truncate(24 * time.Hour).UTC()

		results = append(results, todo)
	}

	return results, err
}

func (s *TodoStore) CreateTodo(t *todoPchi.Todo) error {
	cursor, err := s.InsertOne(context.TODO(), bson.D{
		{"title", t.Title},
		{"description", t.Description},
		{"date", t.Date},
		{"done", false},
	})

	if err != nil {
		return err
	}
	log.Printf("inserted document %v\n", cursor.InsertedID)
	return nil
}

func (s *TodoStore) TodoByID(id primitive.ObjectID) (*todoPchi.Todo, error) {
	var result todoPchi.Todo
	filter := bson.D{{"_id", id}}

	err := s.FindOne(context.TODO(), filter).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) { // errors.Is == (err == mongo.ErrNoDocuments) <-- IDE Told me that
		return nil, fmt.Errorf("TODO with ID %s not found", id)
	} else if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &result, nil
}

func (s *TodoStore) DeleteTodo(id primitive.ObjectID) error {
	filter := bson.D{{"_id", id}}
	cursor, err := s.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Print(err)
		return err
	}
	fmt.Printf("Deleted %v documents with the title '%s'\n", cursor.DeletedCount, id)
	return nil
}

func (s *TodoStore) UpdateTodo(id primitive.ObjectID) error {
	filter := bson.D{{"_id", id}}

	// Поиск документа
	result := s.FindOne(context.TODO(), filter)
	if result.Err() == mongo.ErrNoDocuments {
		log.Println("Document not found")
		// Обработка отсутствия документа
		return nil
	} else if result.Err() != nil {
		log.Fatal(result.Err())
		return result.Err()
	}

	// Декодирование документа
	var todo todoPchi.Todo
	err := result.Decode(&todo)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// Инвертирование значения поля done
	todo.Done = !todo.Done

	// Обновление документа
	update := bson.D{{"$set", bson.D{{"done", todo.Done}}}}
	_, err = s.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Printf("Updated Todo: %+v\n", todo)
	return nil
}

func (s *TodoStore) CreateComment(id primitive.ObjectID) error {
	filter := bson.D{{"id", id}}

	result := s.FindOne(context.TODO(), filter)
	if result.Err() == mongo.ErrNoDocuments {
		log.Println("Document not found")
		// Обработка отсутствия документа
		return nil
	} else if result.Err() != nil {
		log.Fatal(result.Err())
		return result.Err()
	}
	var todo todoPchi.Todo
	err := result.Decode(&todo)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (s *TodoStore) AddComment(id primitive.ObjectID, comment todoPchi.Comment) error {
	filter := bson.D{{"_id", id}}
	update := bson.D{
		{"$push", bson.D{
			{"comments", comment},
		}},
	}

	result := s.FindOneAndUpdate(context.TODO(), filter, update)
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}
