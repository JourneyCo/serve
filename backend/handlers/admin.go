package handlers

import (
        "database/sql"
        "encoding/json"
        "net/http"
        "strconv"
        "time"

        "github.com/gorilla/mux"

        "project-registration-system/middleware"
        "project-registration-system/models"
)

// AdminHandler handles admin-related requests
type AdminHandler struct {
        DB *sql.DB
}

// ProjectInput represents the input for creating or updating a project
type ProjectInput struct {
        Title       string  `json:"title"`
        Description string  `json:"description"`
        StartTime   string  `json:"startTime"`
        EndTime     string  `json:"endTime"`
        ProjectDate string  `json:"projectDate"`
        MaxCapacity int     `json:"maxCapacity"`
        LocationName string  `json:"locationName"`
        Latitude    float64 `json:"latitude"`
        Longitude   float64 `json:"longitude"`
}

// RegisterAdminRoutes registers the routes for admin handlers
func RegisterAdminRoutes(router *mux.Router, db *sql.DB) {
        handler := &AdminHandler{
                DB: db,
        }

        router.HandleFunc("/users", handler.GetAllUsers).Methods("GET")
        router.HandleFunc("/users/{id}/toggle-admin", handler.ToggleUserAdmin).Methods("POST")
        router.HandleFunc("/projects", handler.CreateProject).Methods("POST")
        router.HandleFunc("/projects/{id:[0-9]+}", handler.UpdateProject).Methods("PUT")
        router.HandleFunc("/projects/{id:[0-9]+}", handler.DeleteProject).Methods("DELETE")
}

// GetAllUsers returns all users
func (h *AdminHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
        users, err := models.GetAllUsers(h.DB)
        if err != nil {
                middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve users")
                return
        }

        middleware.RespondWithJSON(w, http.StatusOK, users)
}

// ToggleUserAdmin toggles the admin status of a user
func (h *AdminHandler) ToggleUserAdmin(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        userID := vars["id"]

        // Get the user
        user, err := models.GetUserByID(h.DB, userID)
        if err != nil {
                middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve user")
                return
        }

        if user == nil {
                middleware.RespondWithError(w, http.StatusNotFound, "User not found")
                return
        }

        // Toggle admin status
        err = models.SetUserAdmin(h.DB, userID, !user.IsAdmin)
        if err != nil {
                middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to update user admin status")
                return
        }

        user.IsAdmin = !user.IsAdmin
        middleware.RespondWithJSON(w, http.StatusOK, user)
}

// CreateProject creates a new project
func (h *AdminHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
        var input ProjectInput
        if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
                middleware.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
                return
        }

        // Validate input
        if input.Title == "" || input.Description == "" || input.StartTime == "" || input.EndTime == "" || input.ProjectDate == "" || input.MaxCapacity <= 0 {
                middleware.RespondWithError(w, http.StatusBadRequest, "All fields are required and max capacity must be greater than 0")
                return
        }

        // Parse project date
        projectDate, err := time.Parse("2006-01-02", input.ProjectDate)
        if err != nil {
                middleware.RespondWithError(w, http.StatusBadRequest, "Invalid project date format (use YYYY-MM-DD)")
                return
        }

        // Create project
        project := &models.Project{
                Title:       input.Title,
                Description: input.Description,
                StartTime:   input.StartTime,
                EndTime:     input.EndTime,
                ProjectDate: projectDate,
                MaxCapacity: input.MaxCapacity,
                LocationName: input.LocationName,
                Latitude:    input.Latitude,
                Longitude:   input.Longitude,
        }

        if err := models.CreateProject(h.DB, project); err != nil {
                middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to create project")
                return
        }

        middleware.RespondWithJSON(w, http.StatusCreated, project)
}

// UpdateProject updates an existing project
func (h *AdminHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id, err := strconv.Atoi(vars["id"])
        if err != nil {
                middleware.RespondWithError(w, http.StatusBadRequest, "Invalid project ID")
                return
        }

        // Check if project exists
        project, err := models.GetProjectByID(h.DB, id)
        if err != nil {
                middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve project")
                return
        }

        if project == nil {
                middleware.RespondWithError(w, http.StatusNotFound, "Project not found")
                return
        }

        var input ProjectInput
        if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
                middleware.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
                return
        }

        // Validate input
        if input.Title == "" || input.Description == "" || input.StartTime == "" || input.EndTime == "" || input.ProjectDate == "" || input.MaxCapacity <= 0 {
                middleware.RespondWithError(w, http.StatusBadRequest, "All fields are required and max capacity must be greater than 0")
                return
        }

        // Parse project date
        projectDate, err := time.Parse("2006-01-02", input.ProjectDate)
        if err != nil {
                middleware.RespondWithError(w, http.StatusBadRequest, "Invalid project date format (use YYYY-MM-DD)")
                return
        }

        // Update project
        project.Title = input.Title
        project.Description = input.Description
        project.StartTime = input.StartTime
        project.EndTime = input.EndTime
        project.ProjectDate = projectDate
        project.MaxCapacity = input.MaxCapacity
        project.LocationName = input.LocationName
        project.Latitude = input.Latitude
        project.Longitude = input.Longitude

        if err := models.UpdateProject(h.DB, project); err != nil {
                middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to update project")
                return
        }

        middleware.RespondWithJSON(w, http.StatusOK, project)
}

// DeleteProject deletes a project
func (h *AdminHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id, err := strconv.Atoi(vars["id"])
        if err != nil {
                middleware.RespondWithError(w, http.StatusBadRequest, "Invalid project ID")
                return
        }

        // Check if project exists
        project, err := models.GetProjectByID(h.DB, id)
        if err != nil {
                middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve project")
                return
        }

        if project == nil {
                middleware.RespondWithError(w, http.StatusNotFound, "Project not found")
                return
        }

        // Delete project
        if err := models.DeleteProject(h.DB, id); err != nil {
                middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to delete project")
                return
        }

        middleware.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Project deleted successfully"})
}
