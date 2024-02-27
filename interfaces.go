package todoPchi

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Todo struct {
	ID          primitive.ObjectID `bson:"_id"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
	Date        time.Time          `bson:"date"`
	Done        bool               `bson:"done"`
	Comments    []Comment          `bson:"comments"`
}

type Comment struct {
	ID          primitive.ObjectID `bson:"_id"`
	Description string             `bson:"description"`
	Date        time.Time          `bson:"date"`
}

type TodoStore interface {
	Todos() ([]Todo, error)
	CreateTodo(t *Todo) error
	TodoByID(id primitive.ObjectID) (*Todo, error)
	DeleteTodo(id primitive.ObjectID) error
	UpdateTodo(id primitive.ObjectID) error
}

type CommentStore interface {
	Comments(id primitive.ObjectID) ([]Comment, error)
	CommentCreate(id primitive.ObjectID, t *Comment) error
	AddComment(id primitive.ObjectID, comment Comment) error
}

type Store interface {
	TodoStore
	CommentStore
}
