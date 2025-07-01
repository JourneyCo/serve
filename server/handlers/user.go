package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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
	router.HandleFunc("/registrations/{id:[0-9]+}", handler.UpdateRegistrationGuestCount).Methods(http.MethodPut)
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
			// TODO: check for email address - now that we're down to google and apple, there should
			// be an email ALWAYS
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
	params := r.URL.Query()
	email := params.Get("email")

	if email == "" {
		middleware.RespondWithError(w, http.StatusPreconditionRequired, "Failed to get user information")
	}

	// Get user ID from the token
	user, err := models.GetUserByEmail(ctx, h.DB, email)
	if err != nil || user == nil {
		middleware.RespondWithError(w, http.StatusUnauthorized, "Failed to get user information")
		return
	}

	registration, err := models.GetUserRegistration(ctx, h.DB, user.ID)
	if err != nil {
		log.Println("failed to retrieve user registrations")
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve registrations")
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, registration)
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

// UpdateRegistrationGuestCount updates the guest count for a registration
func (h *UserHandler) UpdateRegistrationGuestCount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	regID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println("invalid registration id to update registration")
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid registration ID")
		return
	}

	var input struct {
		GuestCount int `json:"guest_count"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Println("invalid payload to update registration")
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.GuestCount < 0 {
		log.Println("invalid guest count to update registration")
		middleware.RespondWithError(w, http.StatusBadRequest, "Guest count cannot be negative")
		return
	}

	query := `UPDATE registrations SET guest_count = $1 WHERE id = $2`
	result, err := h.DB.Exec(query, input.GuestCount, regID)
	if err != nil {
		log.Println("failed to update guest count: ", err)
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to update registration")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		log.Println("failed to find registration for updating guest count: ", err)
		middleware.RespondWithError(w, http.StatusNotFound, "Registration not found")
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Registration updated successfully"})
}
