package handlers

import (
	"encoding/json"
	"net/http"

	"serve/middleware"
	"serve/services"
)

// GeocodingHandler handles geocoding requests
type GeocodingHandler struct {
	MapsService *services.MapsService
}

// GeocodeRequest represents a request to geocode an address
type GeocodeRequest struct {
	Address string `json:"address"`
}

// GeocodeAddress handles requests to geocode an address
func (h *GeocodingHandler) GeocodeAddress(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var req GeocodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Geocode the address
	result, err := h.MapsService.GeocodeAddress(req.Address)
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to geocode address: "+err.Error())
		return
	}

	// Return the result
	middleware.RespondWithJSON(w, http.StatusOK, result)
}
