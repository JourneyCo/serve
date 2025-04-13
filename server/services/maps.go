package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

// MapsService handles interactions with the Google Maps API
type MapsService struct {
	apiKey string
}

// GeocodingResult represents the result from geocoding
type GeocodingResult struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"formatted_address"`
}

// GooglePlacesResponse represents the response from Google Places API
type GooglePlacesResponse struct {
	Results []struct {
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"results"`
	Status string `json:"status"`
}

// NewMapsService creates a new Maps service
func NewMapsService() *MapsService {
	return &MapsService{
		apiKey: os.Getenv("GOOGLE_MAPS_API_KEY"),
	}
}

// GeocodeAddress converts an address to latitude and longitude
func (s *MapsService) GeocodeAddress(address string) (*GeocodingResult, error) {
	// Build the request URL
	baseURL := "https://maps.googleapis.com/maps/api/place/textsearch/json"
	params := url.Values{}
	params.Add("query", address)
	params.Add("key", s.apiKey)
	requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Send the request
	resp, err := http.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("error sending geocoding request: %v", err)
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading geocoding response: %v", err)
	}

	var placesResponse GooglePlacesResponse
	if err := json.Unmarshal(body, &placesResponse); err != nil {
		return nil, fmt.Errorf("error parsing geocoding response: %v", err)
	}

	// Check if the response has results
	if placesResponse.Status != "OK" || len(placesResponse.Results) == 0 {
		return nil, fmt.Errorf("no results found for address: %s", address)
	}

	// Get the first result
	result := placesResponse.Results[0]
	return &GeocodingResult{
		Latitude:  result.Geometry.Location.Lat,
		Longitude: result.Geometry.Location.Lng,
		Address:   result.FormattedAddress,
	}, nil
}
