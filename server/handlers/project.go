package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"serve/middleware"
	"serve/models"
	"serve/services"
)

// ProjectHandler handles project-related requests
type ProjectHandler struct {
	DB           *sql.DB
	EmailService *services.EmailService
}

// regRequest defines the JSON request for registration
type regRequest struct {
	GuestCount       int    `json:"guest_count"`
	IsLeadInterested bool   `json:"lead_interest"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Phone            string `json:"phone"`
	Email            string `json:"email"`
	TextPerm         bool   `json:"text_permission"`
}

// RegisterProjectRoutes registers the routes for project handlers
func RegisterProjectRoutes(router *mux.Router, db *sql.DB, emailService *services.EmailService) {
	handler := &ProjectHandler{
		DB:           db,
		EmailService: emailService,
	}

	router.HandleFunc("", handler.GetProjects).Methods("GET")
	router.HandleFunc("/{id:[0-9]+}", handler.GetProject).Methods("GET")
	router.HandleFunc("/{id:[0-9]+}/register", handler.RegisterForProject).Methods("POST")
	router.HandleFunc("/{id:[0-9]+}/cancel", handler.CancelRegistration).Methods("POST")
	router.HandleFunc("/{id:[0-9]+}/registrations", handler.GetProjectRegistrations).Methods("GET")
}

// GetProjects returns all projects
func (h *ProjectHandler) GetProjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projects, err := models.GetAllProjects(ctx, h.DB)
	if err != nil {
		log.Println("error getting projects: ", err)
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve projects")
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, projects)
}

// GetProject returns a specific project by ID
func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	project, err := models.GetProjectByID(ctx, h.DB, id)
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve project")
		return
	}

	if project == nil {
		middleware.RespondWithError(w, http.StatusNotFound, "Project not found")
		return
	}

	// Get serve lead details if serve lead ID exists
	if project.ServeLeadID != "" {
		serveLead, err := models.GetUserByID(ctx, h.DB, project.ServeLeadID)
		if err == nil && serveLead != nil {
			project.ServeLead = serveLead
		}
	}

	middleware.RespondWithJSON(w, http.StatusOK, project)
}

// RegisterForProject registers a user for a project
func (h *ProjectHandler) RegisterForProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	var reg regRequest
	// Parse registration request
	if err = middleware.ParseJSON(r, &reg); err != nil {
		// If parsing fails, use default values (for backward compatibility)
		reg.GuestCount = 0
		reg.IsLeadInterested = false
	}
	err = nil

	// // Get user ID from the token
	// userID, err := middleware.GetUserIDFromRequest(r)
	// if err != nil {
	// 	middleware.RespondWithError(w, http.StatusUnauthorized, "Failed to get user information")
	// 	return
	// }

	// Get or create user in database
	user, err := models.GetUserByEmail(ctx, h.DB, reg.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to check user status")
		return
	}

	var id string
	if user != nil {
		id = user.ID
	}

	// If user doesn't exist, create them
	if user == nil || errors.Is(err, sql.ErrNoRows) {
		var err error
		id, err := uuid.NewUUID()
		if err != nil {
			middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to check user status")
			return
		}
		user = &models.User{
			ID:             id.String(),
			FirstName:      reg.FirstName,
			LastName:       reg.LastName,
			Phone:          reg.Phone,
			Email:          reg.Email,
			TextPermission: reg.TextPerm,
		}
		err = models.CreateUser(ctx, h.DB, user)
		if err != nil {
			middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to create user")
			return
		}
	}

	// Update user information
	//
	// if err = models.UpdateUser(ctx, h.DB, user); err != nil {
	// 	log.Printf("Failed to update user information: %v\n", err)
	// 	middleware.RespondWithError(
	// 		w, http.StatusInternalServerError, "Failed to update user information",
	// 	)
	// 	return
	// }

	// Validate guest count
	if reg.GuestCount < 0 {
		middleware.RespondWithError(w, http.StatusBadRequest, "Guest count cannot be negative")
		return
	}

	// Register for the project
	registration, err := models.RegisterForProject(
		h.DB, id, projectID, reg.GuestCount, reg.IsLeadInterested,
	)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get project details for email
	project, err := models.GetProjectByID(ctx, h.DB, projectID)
	if err != nil {
		middleware.RespondWithError(
			w, http.StatusInternalServerError, "Registration successful but failed to send confirmation email",
		)
		return
	}

	// Get user details for email
	user, err = models.GetUserByID(ctx, h.DB, id)
	if err != nil {
		middleware.RespondWithError(
			w, http.StatusInternalServerError, "Registration successful but failed to send confirmation email",
		)
		return
	}

	// Send confirmation email
	if user != nil && project != nil {
		go func() {
			err = h.EmailService.SendRegistrationConfirmation(user, project)
			if err != nil {
				log.Println("error sending registration email to: ", user.FirstName, " ", user.LastName)
				log.Println(err)
			}
		}()
	}

	middleware.RespondWithJSON(w, http.StatusCreated, registration)
}

// CancelRegistration cancels a user's registration for a project
func (h *ProjectHandler) CancelRegistration(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	// Get user ID from the token
	userID, err := middleware.GetUserIDFromRequest(r)
	if err != nil {
		middleware.RespondWithError(w, http.StatusUnauthorized, "Failed to get user information")
		return
	}

	// Cancel the registration
	err = models.CancelRegistration(h.DB, userID, projectID)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Registration cancelled successfully"})
}

// GetProjectRegistrations returns all registrations for a project (admin only)
func (h *ProjectHandler) GetProjectRegistrations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	registrations, err := models.GetProjectRegistrations(ctx, h.DB, projectID)
	if err != nil {
		log.Println("failed to get project registrations")
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve registrations")
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, registrations)
}
