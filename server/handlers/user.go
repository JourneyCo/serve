package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"project-registration-system/middleware"
	"project-registration-system/models"
	"project-registration-system/services"
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
	router.HandleFunc("/registrations", handler.GetUserRegistrations).Methods("GET")
}

// GetUserProfile returns the profile of the authenticated user
func (h *UserHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
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
	user, err := models.GetUserByID(h.DB, userID)
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
		err = models.CreateUser(h.DB, user)
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
	// Get user ID from the token
	userID, err := middleware.GetUserIDFromRequest(r)
	if err != nil {
		middleware.RespondWithError(w, http.StatusUnauthorized, "Failed to get user information")
		return
	}

	registrations, err := models.GetUserRegistrations(h.DB, userID)
	if err != nil {
		log.Println("failed to retrieve user registrations")
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve registrations")
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, registrations)
}
