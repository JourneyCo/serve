package project

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	db "serve/app/database"
	"strconv"
)

// toCtx will place the project into the context.
func toCtx(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			log.Println("error in parsing project id: ", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		p, err := db.GetProject(ctx, id)
		if err != nil {
			log.Println("error getting project from db: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx = context.WithValue(ctx, "project", p)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// register will register a number of members to a project.
func register(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
