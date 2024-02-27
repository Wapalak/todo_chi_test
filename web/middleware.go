package web

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"net/http"
)

func UUIDFromURLParam(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		todoID := chi.URLParam(r, "_id") // Use the correct parameter name
		uuid, err := uuid.Parse(todoID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Place UUID in the request context
		ctx := context.WithValue(r.Context(), "todoID", uuid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
