package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var TotalProjectsOnSheet = 53
var serveDay = time.Date(2025, 7, 12, 0, 0, 0, 0, time.UTC)
var serveDayPostgresStyle = serveDay.Format("2006-01-02 15:04:05-07:00") // "2025-07-12 00:00:00+00:00"

func main() {
	typeMap1 := make(map[string]int)
	var typInserts []typeInsert
	var projects []Project
	file, err := os.Open("csv/data.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Customize the reader if needed
	// reader.Comma = ';'
	// reader.Comment = '#'

	// Read all records at once
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	typeMapID := 1 // holds the ID number for the db insert of a type
	for line, record := range records {
		// actual projects start at line 2 of the csv
		if line <= 1 {
			continue
		}
		projectID := line - 1

		if line == TotalProjectsOnSheet+1 {
			break
		}

		gID, err := strconv.Atoi(record[0])
		if err != nil {
			log.Fatalf("record number %d has invalid project ID number - %s\n", line, record[0])
		}

		need, err := strconv.Atoi(record[5])
		if err != nil {
			log.Fatalf("record number %d has invalid need number - %s\n", line, record[5])

		}

		lat, long := GeocodeAddress(record[13])

		p := Project{
			GoogleID:        gID,
			Title:           sanitize(record[1]),
			Area:            record[3],
			Time:            record[4],
			MaxCapacity:     need,
			Ages:            record[10],
			LocationAddress: record[13],
			Description:     sanitize(record[16]),
			Website:         record[17],
			Latitude:        lat,
			Longitude:       long,
		}

		projects = append(projects, p)

		if record[2] == "" { // check for types; if none present we continue to next record
			continue
		}
		// types present
		types := strings.Split(record[2], "\n")

		for _, t := range types {
			val, ok := typeMap1[t]
			if !ok {
				typeMap1[t] = typeMapID
				val = typeMap1[t]
				typeMapID++
			}
			ci := typeInsert{
				projectID: projectID,
				typeID:    val,
			}
			typInserts = append(typInserts, ci)
		}
	}

	if len(projects) == 0 {
		return
	}

	sqlStmt := `INSERT INTO projects (google_id, title, description, time, project_date, max_capacity, area, latitude, longitude, serve_lead_id, location_address, website, ages) VALUES `

	for _, val := range projects {
		vals := fmt.Sprintf(
			"(%d, '%s', '%s', '%s', '%s', %d, '%s', %f, %f, '%s', '%s', '%s', '%s'), ", val.GoogleID, val.Title,
			val.Description,
			val.Time, serveDayPostgresStyle, val.MaxCapacity, val.Area, val.Latitude, val.Longitude,
			"example-user-123", val.LocationAddress, val.Website, val.Ages,
		)
		sqlStmt += vals
	}

	sqlStmt, ok := strings.CutSuffix(sqlStmt, ", ")

	if !ok {
		log.Fatal("couldn't find the end of sql statement")
	}

	sqlStmt += ";"

	// if types are not included, we are done and can return just a
	// projects insert
	if len(typInserts) == 0 {
		log.Println("\n\n\nThe Final Sql Statement: ")
		log.Print(sqlStmt)
		return
	}

	// build out the types table dynamically from the types in the spreadsheet
	typesCreateStmt := `INSERT INTO types (id, type)  VALUES `
	for typString, typeID := range typeMap1 {
		vals := fmt.Sprintf("(%d, '%s'), ", typeID, typString)
		typesCreateStmt += vals
	}

	typesCreateStmt, ok = strings.CutSuffix(typesCreateStmt, ", ")
	if !ok {
		log.Fatal("couldn't find the end of sql statement")
	}

	typesCreateStmt += ";"

	typesStmt := `INSERT INTO project_types (project_id, type_id)  VALUES `

	for _, val := range typInserts {
		vals := fmt.Sprintf("(%d, '%d'), ", val.projectID, val.typeID)
		typesStmt += vals
	}

	typesStmt, ok = strings.CutSuffix(typesStmt, ", ")
	if !ok {
		log.Fatal("couldn't find the end of sql statement")
	}

	typesStmt += ";"

	log.Println("\n\n\nThe Final Sql Statement: ")
	log.Printf("%s\n%s\n%s", sqlStmt, typesCreateStmt, typesStmt)
}

func sanitize(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

// GeocodeAddress converts an address to latitude and longitude
func GeocodeAddress(address string) (float64, float64) {
	if strings.Contains(strings.ToLower(address), "tbd") || strings.TrimSpace(address) == "" {
		return 39.491482, -104.874878
	}
	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")

	// Build the request URL
	baseURL := "https://maps.googleapis.com/maps/api/place/textsearch/json"
	params := url.Values{}
	params.Add("query", address)
	params.Add("key", apiKey)
	requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Send the request
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Fatal("could not connect to the api")
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("error reading body")
	}

	var placesResponse GooglePlacesResponse
	if err := json.Unmarshal(body, &placesResponse); err != nil {
		log.Fatal("could not parse")
	}

	// Check if the response has results
	if placesResponse.Status != "OK" || len(placesResponse.Results) == 0 {
		log.Print("got here")
		log.Fatal("could not get result")
	}

	// Get the first result
	result := placesResponse.Results[0]
	return result.Geometry.Location.Lat, result.Geometry.Location.Lng
}

type Project struct {
	GoogleID        int       `json:"google_id"`
	Title           string    `json:"title"`       // Project
	Description     string    `json:"description"` // About this project
	Website         string    `json:"website"`
	Time            string    `json:"time"`
	ProjectDate     time.Time `json:"project_date"`
	MaxCapacity     int       `json:"max_capacity"` // Volunteers
	CurrentReg      int       `json:"current_registrations"`
	Area            string    `json:"area"`
	LocationAddress string    `json:"location_address"` // Address
	Latitude        float64   `json:"latitude"`
	Longitude       float64   `json:"longitude"`
	ServeLeadID     string    `json:"serve_lead_id"`
	Types           []string  `json:"types,omitempty"` // Type
	Ages            string    `json:"ages,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type typeInsert struct {
	projectID int
	typeID    int
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
