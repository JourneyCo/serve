package project

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	db "serve/app/database"
	"serve/models"
)

// toCtx will place the project into the context.
func toCtx(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
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
		},
	)
}

// register will register a number of members to a project.
func register(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			proj := ctx.Value("project").(models.Project)
			body := ctx.Value("body").([]byte)
			var dto Request
			now := time.Now()

			if err := json.Unmarshal(body, &dto); err != nil {
				fmt.Printf("error unmarshalling body: %v", err)
				return
			}

			if dto.Registering == nil {
				log.Print("request did not include members to register")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			toRegister := *dto.Registering

			if proj.Needed <= 0 || toRegister > proj.Needed {
				log.Printf("project has a need of %d but user attempted to register %d", proj.Needed, toRegister)
				w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
				return
			}

			proj.Needed -= toRegister
			proj.UpdatedAt = &now

			project, err := db.PutProject(ctx, proj)
			if err != nil {
				log.Printf("failed to put project: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			ctx = context.WithValue(ctx, "project", project)

			h.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}
