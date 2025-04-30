package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Project struct {
	GoogleID         int
	Categories       string
	OrgName          string
	ShortDescription string
	OrgContact       string
	OrgNumber        string
	OrgLead          string
	ProjectType      string
	Location         string
	Time             string
	Needed           int
	Ages             string
	Lead             string
	Skills           string
	Supplies         string
	Tools            string
	LongDescription  string
	OnApp            string
	ProjectQty       int
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

		gID, err := strconv.Atoi(record[0])
		if err != nil {
			log.Printf("record number %d has invalid project ID number - %s\n", line, record[0])
			gID = 0
			err = nil
			continue // TODO: remove if you want all of the records!!!
		}

		need, err := strconv.Atoi(record[10])
		if err != nil {
			log.Printf("record number %d has invalid need number - %s\n", line, record[10])
			need = 0
			err = nil
			continue // TODO: remove if you want all of the records!!!

		}

		pqty, err := strconv.Atoi(record[19])
		if err != nil {
			log.Printf("record number %d has invalid project qty number - %s\n", line, record[18])
			pqty = 0
			err = nil
			continue // TODO: remove if you want all of the records!!!

		}

		p := Project{
			GoogleID:         gID,
			Categories:       record[1],
			OrgName:          record[2],
			ShortDescription: record[3],
			OrgContact:       record[4],
			OrgNumber:        record[5],
			OrgLead:          record[6],
			ProjectType:      record[7],
			Location:         record[8],
			Time:             record[9],
			Needed:           need,
			Ages:             record[11],
			Lead:             record[13],
			Skills:           record[14],
			Supplies:         record[15],
			Tools:            record[16],
			LongDescription:  record[17],
			OnApp:            record[18],
			ProjectQty:       pqty,
		}

		projects = append(projects, p)
	}

	if len(projects) == 0 {
		return
	}

	sqlStmt := `INSERT INTO projects (google_id, title, short_description, description, time, project_date,
		max_capacity, location_name, latitude, longitude, lead_user_id, wheelchair_accessible, location_address
	) VALUES `

	for _, val := range projects {
		vals := fmt.Sprintf(
			"(%d, '%s', '%s', '%s', '%s', '%s', %d, '%s', %f, %f, '%s', %t, '%s'), ", val.GoogleID, val.OrgName,
			val.ShortDescription,
			val.LongDescription, val.Time, "2025-07-12", val.Needed, val.Location, 39.491482, -104.874878,
			"example-user-123", true, val.Location,
		)
		sqlStmt += vals
	}

	sqlStmt, ok := strings.CutSuffix(sqlStmt, ", ")
	if !ok {
		log.Fatal("couldn't find the end of sql statement")
	}

	sqlStmt += ";"

	log.Println("\n\n\nThe Final Sql Statement: \n")
	log.Print(sqlStmt)

}
