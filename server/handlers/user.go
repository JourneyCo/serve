package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"serve/middleware"
	"serve/models"
	"serve/services"
)

// UserHandler handles user-related requests
type UserHandler struct {
	DB           *sql.DB
	EmailService *services.EmailService
}

// RegisterUserRoutes registers the routes for user handlers
func RegisterUserRoutes(router *mux.Router, db *sql.DB, emailService *services.EmailService) {
	handler := &UserHandler{
		DB:           db,
		EmailService: emailService,
	}

	router.HandleFunc("/profile", handler.GetUserProfile).Methods("GET")
	router.HandleFunc("/profile", handler.UpdateUserProfile).Methods("PUT")
	router.HandleFunc("/registrations", handler.GetUserRegistrations).Methods("GET")
}

// GetUserProfile returns the profile of the authenticated user
func (h *UserHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user ID from the token
	userID, err := middleware.GetUserIDFromRequest(r)
	if err != nil {
		middleware.RespondWithError(w, http.StatusUnauthorized, "Failed to get user information")
		return
	}

	// Get user claims from token
	_, err = middleware.GetUserFromRequest(r)
	if err != nil {
		middleware.RespondWithError(w, http.StatusUnauthorized, "Failed to get user information")
		return
	}

	// Get or create user in database
	user, err := models.GetUserByID(ctx, h.DB, userID)
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve user profile")
		return
	}

	// If user doesn't exist, create a new one
	if user == nil {
		user = &models.User{
			ID: userID,
		}

		// Create the user in the database
		err = models.CreateUser(ctx, h.DB, user)
		if err != nil {
			middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to create user profile")
			return
		}
	} else {
		// Update user information if it has changed
		// if user.Email != claims.Email || user.Name != claims.Name || user.Picture != claims.Picture {
		//         user.Email = claims.Email
		//         user.Name = claims.Name
		//         user.Picture = claims.Picture
		//
		//         err = models.UpdateUser(h.DB, user)
		//         if err != nil {
		//                 middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to update user profile")
		//                 return
		//         }
		// }
		// TODO: update this
	}

	middleware.RespondWithJSON(w, http.StatusOK, user)
}

// GetUserRegistrations returns all registrations for the authenticated user
func (h *UserHandler) GetUserRegistrations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user ID from the token
	userID, err := middleware.GetUserIDFromRequest(r)
	if err != nil {
		middleware.RespondWithError(w, http.StatusUnauthorized, "Failed to get user information")
		return
	}

	registrations, err := models.GetUserRegistrations(ctx, h.DB, userID)
	if err != nil {
		log.Println("failed to retrieve user registrations")
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve registrations")
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, registrations)
}

// UpdateUserProfile updates the profile of the authenticated user
func (h *UserHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := middleware.GetUserIDFromRequest(r)
	if err != nil {
		middleware.RespondWithError(w, http.StatusUnauthorized, "Failed to get user information")
		return
	}

	var user models.User
	if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Ensure the user can only update their own profile
	user.ID = userID

	if err = models.UpdateUser(ctx, h.DB, &user); err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to update user profile")
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, user)
}
