package web

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"html/template"
	"log"
	"net/http"
	"time"
	"todoPchi"
	"todoPchi/database"
)

type Handler struct {
	*chi.Mux
	todo todoPchi.Store
}

func NewHadler(t *database.Store) *Handler {
	h := &Handler{
		Mux:  chi.NewMux(),
		todo: t,
	}

	h.Use(middleware.Logger)
	h.Route("/todos", func(r chi.Router) {
		r.Get("/", h.TodoList())

		r.Get("/new", h.TodoCreate())
		r.Post("/", h.TodoStore())
		r.Post("/{id}/delete", h.TodoDelete())
		r.Get("/{id}/todo", h.TodoView())
		r.Post("/{id}/todo/update", h.TodoUpdate())
		r.Get("/{id}/todo/AddComment", h.CommentCreate())
		r.Post("/{id}/todo/AddComment/Create", h.CommentPost())

		//r.Post("/{id}/todo", h.CommentAdd())
		//r.Get("/{id}/todo/AddComment", h.CommentAdd())    // Новый роут для страницы добавления комментария
		//r.Post("/{id}/todo/AddComment", h.CommentStore()) // Новый роут для обработки добавления комментария
	})

	return h
}

func (h *Handler) TodoList() http.HandlerFunc {
	type data struct {
		Todos []todoPchi.Todo
	}
	tmpl := template.Must(template.New("todoList.html").ParseFiles("C:\\Users\\User\\GolandProjects" +
		"\\todoPchi\\web\\templates\\todoList.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		tt, err := h.todo.Todos()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, data{Todos: tt})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handler) TodoCreate() http.HandlerFunc {
	tmpl := template.Must(template.New("todoCreate.html").ParseFiles("C:\\Users\\User\\GolandProjects" +
		"\\todoPchi\\web\\templates\\todoCreate.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.Execute(w, nil); err != nil {
			log.Fatal(err)
		}
	}
}

func (h *Handler) TodoStore() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.FormValue("title")
		description := r.FormValue("description")
		dateStr := r.FormValue("date")

		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			http.Error(w, "Invalid date format", http.StatusBadRequest)
			return
		}

		if err := h.todo.CreateTodo(&todoPchi.Todo{
			ID:          primitive.NewObjectID(),
			Title:       title,
			Description: description,
			Date:        date,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/todos", http.StatusFound)
	}
}

func (h *Handler) TodoView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")

		id, err := primitive.ObjectIDFromHex(idStr)

		todo, err := h.todo.TodoByID(id)
		log.Println(id)
		comment, err := h.todo.Comments(id)

		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				http.Error(w, "TODO not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
			}
			return
		}
		data := struct {
			Todo    todoPchi.Todo
			Comment []todoPchi.Comment
		}{
			Todo:    *todo,
			Comment: comment,
		}
		tmpl := template.Must(template.New("todoInfo.html").ParseFiles("C:\\Users\\User\\" +
			"GolandProjects\\todoPchi\\web\\templates\\todoInfo.html"))
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handler) TodoDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")

		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := h.todo.DeleteTodo(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/todos", http.StatusFound)
	}
}

func (h *Handler) TodoUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")

		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if err := h.todo.UpdateTodo(id); err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/todos/"+idStr+"/todo", http.StatusFound)
	}
}

func (h *Handler) CommentCreate() http.HandlerFunc {
	tmpl := template.Must(template.New("commentAdd.html").ParseFiles("C:\\Users\\User\\GolandProjects" +
		"\\todoPchi\\web\\templates\\commentAdd.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Pass Todo ID to the template
		if err := tmpl.Execute(w, id); err != nil {
			log.Fatal(err)
		}
	}
}

func (h *Handler) CommentPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract Todo ID from the URL
		idStr := chi.URLParam(r, "id")
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Parse the form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Get the description from the form
		description := r.FormValue("description")

		// Create a new comment
		comment := todoPchi.Comment{
			ID:          primitive.NewObjectID(),
			Description: description,
			Date:        time.Now(),
		}

		// Update the Todo with the new comment
		if err := h.todo.AddComment(id, comment); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect to the Todo details page
		http.Redirect(w, r, fmt.Sprintf("/todos/%s/todo", idStr), http.StatusFound)
	}
}

//
//func (h *Handler) CommentAdd() http.HandlerFunc {
//	type data struct {
//		Todos []todoPchi.Todo
//	}
//	tmpl := template.Must(template.New("commentAdd.html").ParseFiles(
//		"C:\\Users\\User\\GolandProjects\\todoPchi\\web\\templates\\commentAdd.html"))
//	return func(w http.ResponseWriter, r *http.Request) {
//		idStr := chi.URLParam(r, "id")
//
//		id, err := primitive.ObjectIDFromHex(idStr)
//		if err != nil {
//			log.Println(err)
//			http.Error(w, err.Error(), http.StatusBadRequest)
//			return
//		}
//		todo, err := GetTodoByID(r.Context(), id, h.) // Replace with your actual MongoDB collection
//		if err != nil {
//			log.Println(err)
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//
//		if err := tmpl.Execute(w, data{Todos: todo}); err != nil {
//			http.Error(w, err.Error(), http.StatusBadRequest)
//		}
//	}
//}
//
//func (h *Handler) CommentStore() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		//title := r.FormValue("title")
//		//description := r.FormValue("description")
//
//		err := r.ParseForm()
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusBadRequest)
//			return
//		}
//		//idStr:=  todoHandler(w, "id")
//		idStr := chi.URLParam(r, "id")
//
//		id, err := primitive.ObjectIDFromHex(idStr)
//		if err != nil {
//			log.Println(err)
//			http.Error(w, err.Error(), http.StatusBadRequest)
//			return
//		}
//		//title := r.Form.Get("title")
//		description := r.Form.Get("description")
//		date := time.Now()
//
//		if err := h.todo.CommentCreate(id, &todoPchi.Comment{
//			ID: primitive.NewObjectID(),
//			//TodoID:      id,
//			//Title:       title,
//			Description: description,
//			Date:        date,
//		}); err != nil {
//			log.Println(err)
//			return
//		}
//
//		// Используйте абсолютный путь при редиректе
//		http.Redirect(w, r, "/todos", http.StatusFound)
//	}
//}
//
//func GetTodoByID(ctx context.Context, id primitive.ObjectID, collection *mongo.Collection) (*todoPchi.Todo, error) {
//	var todo todoPchi.Todo
//	filter := bson.M{"_id": id}
//	err := collection.FindOne(ctx, filter).Decode(&todo)
//	if err != nil {
//		return nil, err
//	}
//	return &todo, nil
//}
