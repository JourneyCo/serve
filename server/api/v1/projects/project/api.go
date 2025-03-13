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
	"serve/app/auth0"
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
			session := ctx.Value("session").(auth0.Session)

			if err := json.Unmarshal(body, &dto); err != nil {
				fmt.Printf("error unmarshalling body: %v", err)
				return
			}

			if session.UserID == "" {
				log.Print("session does not include user id")
				w.WriteHeader(http.StatusBadRequest)
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

			var ld bool
			if dto.Lead != nil && *dto.Lead == true {
				ld = true
			}

			reg := models.Registration{
				AccountID:   session.UserID,
				ProjectID:   proj.ID,
				UpdatedAt:   &now,
				QtyEnrolled: toRegister,
				Lead:        &ld,
			}

			registration, err := db.PutRegistration(ctx, reg)
			if err != nil {
				log.Printf("failed to put registration: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			log.Printf("user %s registered for project %s", reg.AccountID, proj.Name)

			ctx = context.WithValue(ctx, "project", project)
			ctx = context.WithValue(ctx, "registration", registration)

			h.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}
