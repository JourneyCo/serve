package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var TotalProjectsOnSheet = 25
var serveDay = time.Date(2025, 7, 12, 0, 0, 0, 0, time.UTC)
var serveDayPostgresStyle = serveDay.Format("2006-01-02 15:04:05-07:00") // "2025-07-12 00:00:00+00:00"

type Project struct {
	GoogleID             int       `json:"google_id"`
	Title                string    `json:"title"` // Project
	ShortDescription     string    `json:"short_description"`
	Description          string    `json:"description"` // About this project
	Website              string    `json:"website"`
	Time                 string    `json:"time"`
	ProjectDate          time.Time `json:"project_date"`
	MaxCapacity          int       `json:"max_capacity"` // Volunteers
	CurrentReg           int       `json:"current_registrations"`
	Area                 string    `json:"area"`
	LocationAddress      string    `json:"location_address"` // Address
	Latitude             float64   `json:"latitude"`
	Longitude            float64   `json:"longitude"`
	WheelchairAccessible bool      `json:"wheelchair_accessible"`
	ServeLeadID          string    `json:"serve_lead_id"`
	Tools                []string  `json:"tools,omitempty"`
	Supplies             []string  `json:"supplies,omitempty"`
	Categories           []string  `json:"categories,omitempty"` // Type
	Ages                 []string  `json:"ages,omitempty"`
	Skills               []string  `json:"skills,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

func main() {
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

	var projects []Project

	for line, record := range records {
		if line <= 1 {
			continue
		}
		if line == TotalProjectsOnSheet+1 {
			break
		}

		gID, err := strconv.Atoi(record[0])
		if err != nil {
			log.Printf("record number %d has invalid project ID number - %s\n", line, record[0])
			gID = 0
			err = nil
			continue // TODO: remove if you want all of the records!!!
		}

		need, err := strconv.Atoi(record[5])
		if err != nil {
			log.Printf("record number %d has invalid need number - %s\n", line, record[5])
			need = 0
			err = nil
			continue // TODO: remove if you want all of the records!!!

		}

		p := Project{
			GoogleID:        gID,
			Title:           record[1],
			Area:            record[3],
			Time:            record[4],
			MaxCapacity:     need,
			LocationAddress: record[13],
			Description:     record[16],
			Website:         record[17],
		}

		projects = append(projects, p)
	}

	if len(projects) == 0 {
		return
	}

	sqlStmt := `INSERT INTO projects (google_id, title, short_description, description, time, project_date,
		max_capacity, area, latitude, longitude, serve_lead_id, wheelchair_accessible, location_address
	) VALUES `

	for _, val := range projects {
		vals := fmt.Sprintf(
			"(%d, '%s', '%s', '%s', '%s', '%s', %d, '%s', %f, %f, '%s', %t, '%s'), ", val.GoogleID, val.Title,
			"", val.Description,
			val.Time, serveDayPostgresStyle, val.MaxCapacity, val.Area, 39.491482, -104.874878,
			"example-user-123", true, val.LocationAddress,
		)
		sqlStmt += vals
	}

	sqlStmt, ok := strings.CutSuffix(sqlStmt, ", ")

	if !ok {
		log.Fatal("couldn't find the end of sql statement")
	}

	sqlStmt += ";"

	log.Println("\n\n\nThe Final Sql Statement: ")
	log.Print(sqlStmt)

}
