package handlers

import (
	"database/sql"
	"errors"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"serve/config"
	"serve/middleware"
	"serve/models"
	"serve/services"
)

// ProjectHandler handles project-related requests
type ProjectHandler struct {
	DB           *sql.DB
	EmailService *services.EmailService
	Config       *config.Config
	TextService  *services.TextService
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
	Recaptcha        string `json:"recaptcha"`
}

// RegisterProjectRoutes registers the routes for project handlers
func RegisterProjectRoutes(
	router *mux.Router, db *sql.DB, cfg *config.Config, emailService *services.EmailService,
	textService *services.TextService,
) {
	handler := &ProjectHandler{
		DB:           db,
		EmailService: emailService,
		Config:       cfg,
		TextService:  textService,
	}

	router.HandleFunc("", handler.GetProjects).Methods("GET")
	router.HandleFunc("/my", handler.GetMyProject).Methods("GET")
	router.HandleFunc("/types", handler.GetTypes).Methods("GET")
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

// GetMyProject returns the project for a user that has already signed up
func (h *ProjectHandler) GetMyProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params := r.URL.Query()
	email := params.Get("email")

	if strings.TrimSpace(email) == "" {
		middleware.RespondWithError(w, http.StatusBadRequest, "No email supplied with request")
		return
	}

	user, err := models.GetUserByEmail(ctx, h.DB, email)
	if errors.Is(err, sql.ErrNoRows) {
		log.Println("no registrations found for this email")
		middleware.RespondWithError(w, http.StatusContinue, "Failed to retrieve registrations")
		return
	} else {
		if err != nil {
			middleware.RespondWithError(w, http.StatusTeapot, "Failed to find user")
			return
		}
	}

	registration, err := models.GetUserRegistration(ctx, h.DB, user.ID)
	if errors.Is(err, sql.ErrNoRows) {
		log.Println("no registrations found for this email")
		middleware.RespondWithError(w, http.StatusContinue, "Failed to retrieve registrations")
		return
	} else {
		if err != nil {
			log.Println("failed to retrieve user registrations")
			middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve registrations")
			return
		}
	}

	middleware.RespondWithJSON(w, http.StatusOK, registration)
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

	if strings.TrimSpace(reg.Email) == "" {
		middleware.RespondWithError(w, http.StatusBadRequest, "No email provided with request")
		return
	}

	// check to see if user is already registered with a project
	existProj, err := models.GetUserRegistrationByEmail(ctx, h.DB, reg.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to check user status")
		return
	}

	// case - they are attempting to re-register again. Send a 208 and front end will handle.
	if !errors.Is(err, sql.ErrNoRows) && existProj == projectID {
		middleware.RespondWithError(w, http.StatusAlreadyReported, "Current user is already signed up for this project")
		return
	}

	// case - they are attempting to register for multiple projects. Send a 409 and front end will handle.
	if !errors.Is(err, sql.ErrNoRows) && existProj != math.MaxInt {
		middleware.RespondWithError(w, http.StatusConflict, "Current user is already signed up for a project")
		return
	}

	err = nil

	// if err := services.CreateAssessment(h.Config, reg.Recaptcha); err != nil {
	// 	middleware.RespondWithError(w, http.StatusBadRequest, "Recaptcha validation failed")
	// 	return
	// }

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

	var userID string
	if user != nil {
		userID = user.ID
	}

	// If user doesn't exist, create them
	if user == nil || errors.Is(err, sql.ErrNoRows) {
		err = nil
		uid, err := uuid.NewUUID()
		if err != nil {
			middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to check user status")
			return
		}
		userID = uid.String()
		user = &models.User{
			ID:             userID,
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

	// Validate guest count
	if reg.GuestCount < 0 {
		middleware.RespondWithError(w, http.StatusBadRequest, "Guest count cannot be negative")
		return
	}

	// Register for the project
	registration, err := models.RegisterForProject(
		h.DB, userID, projectID, reg.GuestCount, reg.IsLeadInterested,
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

	// Send confirmation email
	if user != nil && project != nil {
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			err = h.EmailService.SendRegistrationConfirmation(user, project)
			if err != nil {
				log.Println("error sending registration email to: ", user.FirstName, " ", user.LastName)
				log.Println(err)
			}
		}()
		go func() {
			defer wg.Done()
			err = h.TextService.SendRegistrationConfirmation(user, project)
			if err != nil {
				log.Println("error sending registration email to: ", user.FirstName, " ", user.LastName)
				log.Println(err)
			}
		}()

		wg.Wait() // wait for email and text to both send
	}

	middleware.RespondWithJSON(w, http.StatusCreated, registration)
}

// CancelRegistration cancels a user's registration for a project
func (h *ProjectHandler) CancelRegistration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}
	params := r.URL.Query()
	email := params.Get("email")

	// Get user ID from the token
	user, err := models.GetUserByEmail(ctx, h.DB, email)
	if err != nil {
		middleware.RespondWithError(w, http.StatusUnauthorized, "Failed to get user information")
		return
	}

	// Cancel the registration
	err = models.CancelRegistration(ctx, h.DB, user.ID, projectID)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Registration cancelled successfully"})
}

// GetTypes returns all types from the types table
func (h *ProjectHandler) GetTypes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	types, err := models.GetAllTypes(ctx, h.DB)
	if err != nil {
		log.Println("error getting types: ", err)
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve types")
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, types)
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
