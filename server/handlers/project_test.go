
package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetProjects(t *testing.T) {
	// Setup router and handler
	router := mux.NewRouter()
	handler := &ProjectHandler{DB: nil} // Mock DB would be used here
	router.HandleFunc("", handler.GetProjects).Methods("GET")

	// Create test request
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusInternalServerError, w.Code, "Should return 500 when DB is nil")
}

func TestRegisterForProject(t *testing.T) {
	// Setup router and handler
	router := mux.NewRouter()
	handler := &ProjectHandler{DB: nil} // Mock DB would be used here
	router.HandleFunc("/{id:[0-9]+}/register", handler.RegisterForProject).Methods("POST")

	// Create test registration request
	regRequest := regRequest{
		GuestCount: 2,
		IsLeadInterested: true,
		FirstName: "Test",
		LastName: "User",
		Phone: "1234567890",
		Email: "test@example.com",
	}
	body, _ := json.Marshal(regRequest)

	// Create test request
	req := httptest.NewRequest("POST", "/1/register", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusUnauthorized, w.Code, "Should return 401 when no auth token present")
}

func TestCancelRegistration(t *testing.T) {
	// Setup router and handler
	router := mux.NewRouter()
	handler := &ProjectHandler{DB: nil} // Mock DB would be used here
	router.HandleFunc("/{id:[0-9]+}/cancel", handler.CancelRegistration).Methods("POST")

	// Create test request
	req := httptest.NewRequest("POST", "/1/cancel", nil)
	w := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusUnauthorized, w.Code, "Should return 401 when no auth token present")
}
