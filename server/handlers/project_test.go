
package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"github.com/stretchr/testify/assert"
	"serve/testutils"
)

func TestGetProjects(t *testing.T) {
	// Setup test server
	ts := testutils.NewTestServer()
	defer ts.Close()

	// Create test project
	_, err := testutils.CreateTestProject(ts.DB)
	assert.NoError(t, err)

	// Make request to test server
	resp, err := http.Get(ts.Server.URL + "/api/projects")
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Assert response
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Clean up
	err = testutils.CleanTestData(ts.DB)
	assert.NoError(t, err)
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
