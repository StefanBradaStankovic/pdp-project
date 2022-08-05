package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// ---------------------------------------- CONFIG SECTION ----------------------------------------

// Defining a json-friendly struct for rows from 'movie_revenues' table
type MovieRevenues struct {
	RevenueID            *int     `json:"revenueID"`
	MovieID              *int     `json:"movieID"`
	DomesticTakings      *float64 `json:"domesticTakings"`
	InternationalTakings *float64 `json:"internationalTakings"`
}

var revenuesStatement = queryStatement{
	"SELECT revenue_id, movie_id, domestic_takings, international_takings FROM movie_revenues WHERE revenue_id = $1 AND is_visible", //selectSingle
	"SELECT revenue_id, movie_id, domestic_takings, international_takings FROM movie_revenues WHERE is_visible",                     //selectAll
	"INSERT INTO movie_revenues (movie_id, domestic_takings, international_takings) VALUES ($1, $2, $3) RETURNING revenue_id",       //insertInto
	"DELETE FROM movie_revenues WHERE revenue_id = $1 RETURNING revenue_id",                                                         //deleteFrom
	"UPDATE movie_revenues SET is_visible = false WHERE revenues_id = $1",                                                           //updateVisible
	"SELECT count(revenue_id) FROM movie_revenues WHERE revenue_id = $1 AND is_visible"}                                             //checkForID

// ---------------------------------------- SERVICE SECTION ----------------------------------------

// Get a single set of revenues from the database, encode it into a json
// object and send it as a response. Log the activity into console
func getRevenues(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)
	var idParam string = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	// CHECK for a valid ID parameter
	if err != nil {
		w.WriteHeader(statusCodeBadRequest.status)
		w.Write([]byte(statusCodeBadRequest.message + " - ID not an integer!"))
		fmt.Println("ERROR: revenues.go/strconv.Atoi():   Not integer!")
		return
	}
	// CHECK if item of ID exists in the database
	if exists, err := checkIfExists(id, revenuesStatement.checkForID); err != nil || !exists {
		w.WriteHeader(statusCodeNotFound.status)
		return
	}
	// CHECK if there was an error during query
	revenueRow, err := selectItem(id, revenuesStatement.selectSingle)
	if err != nil {
		w.WriteHeader(statusCodeInternalError.status)
		w.Write([]byte(statusCodeInternalError.message + " - could not fetch data!"))
		fmt.Printf("ERROR: revenues.go/selectItem:   %s\n", err)
		return
	}
	// SCAN row data into json-able object
	revenue := revenueRow.ScanRevenues()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(revenue)

	fmt.Printf("SELECT | movie_revenues | %d\n", *revenue.RevenueID)
}

// Get all rows of movie revenues from the database, encode it into a json
// object and send it as a response. Log the activity into console
func getAllRevenues(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)
	var allRevenues []MovieRevenues
	// CHECK if there was an error during query
	revenuesRows, err := selectAllItems(revenuesStatement.selectAll)
	if err != nil {
		w.WriteHeader(statusCodeInternalError.status)
		w.Write([]byte(statusCodeInternalError.message + " - could not fetch data!"))
		fmt.Printf("ERROR: revenues.go/selectItem:   %s\n", err)
		return
	}
	// SCAN rows data into json-able object
	for revenuesRows.rows.Next() {
		revenues := revenuesRows.ScanAllRevenues()
		allRevenues = append(allRevenues, revenues)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allRevenues)

	fmt.Println("SELECT | movie_revenues | ALL")
}

// Create a single set of revenues in a database based on input JSON object
func postRevenues(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)
	var revenues MovieRevenues
	var statement SqlStmt
	var err error
	json.NewDecoder(r.Body).Decode(&revenues)
	// Add error handling
	statement.statement, err = insertItem(revenuesStatement.insertInto)
	if err != nil {
		w.WriteHeader(statusCodeQueryError.status)
		w.Write([]byte(statusCodeQueryError.message))
		fmt.Printf("ERROR: revenues.go/insertItem:   %s\n", err)
		return
	}
	id := statement.InsertIntoRevenues(revenues)
	if id <= 0 {
		w.WriteHeader(statusCodeQueryError.status)
		w.Write([]byte(statusCodeQueryError.message + " - Could not execute query!"))
		fmt.Printf("ERROR: revenues.go/InsertIntoRevenues:   %s\n", err)
		return
	}
	statement.statement.Close()
	revenues.RevenueID = &id

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(revenues)

	fmt.Printf("INSERT INTO | movie_revenues | %d\n", *revenues.RevenueID)
}

// Find a single set of revenues in a database and delete it permanently
func deleteRevenues(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)
	var idParam string = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	// CHECK for a valid ID parameter
	if err != nil {
		w.WriteHeader(statusCodeBadRequest.status)
		w.Write([]byte(statusCodeBadRequest.message + " - ID not an integer!"))
		fmt.Println("ERROR: revenues.go/strconv.Atoi():   Not integer!")
		return
	}
	// CHECK if item of ID exists in the database
	if exists, err := checkIfExists(id, revenuesStatement.checkForID); err != nil || !exists {
		w.WriteHeader(statusCodeNotFound.status)
		w.Write([]byte(statusCodeNotFound.message + " - item does not exist!"))
		return
	}

	// Sending DELETE FROM query
	err = deleteItem(id, revenuesStatement.updateVisible)
	if err != nil {
		fmt.Printf("ERROR: revenues.go/deleteItem: %s\n", err)
		return
	}
	w.WriteHeader(statusCodeItemDeleted.status)
	w.Write([]byte(fmt.Sprintf("%s - %d", statusCodeItemDeleted.message, id)))
	w.Header().Set("Content-Type", "application/json")

	fmt.Printf("DELETE FROM | movie_revenues | %d\n", id)
}

// ---------------------------------------- INTERFACE SECTION ----------------------------------------

// rowScanner method for scanning a 'MovieRevenues' object
func (inputRow *SqlRow) ScanRevenues() MovieRevenues {
	var output MovieRevenues
	err := inputRow.row.Scan(&output.RevenueID, &output.MovieID, &output.DomesticTakings, &output.InternationalTakings)
	if err != nil {
		fmt.Printf("ERROR - revenues.go/interface/ScanRevenues -  %s\n", err)
		return output
	}
	return output
}

// rowScanner method for scanning an 'MovieRevenues' object from multiple rows
func (inputRow *SqlRows) ScanAllRevenues() MovieRevenues {
	var output MovieRevenues
	err := inputRow.rows.Scan(&output.RevenueID, &output.MovieID, &output.DomesticTakings, &output.InternationalTakings)
	if err != nil {
		fmt.Printf("ERROR - directors.go/interface/ScanAllRevenues -  %s\n", err)
		return output
	}
	return output
}

// rowScanner method for inserting an 'MovieRevenues' object into the database as a single row
func (queryStatement *SqlStmt) InsertIntoRevenues(inputRow MovieRevenues) int {
	var output int
	execTime := time.Now().UnixMilli()
	err := queryStatement.statement.QueryRow(inputRow.MovieID, inputRow.DomesticTakings, inputRow.InternationalTakings).Scan(&output)
	if err != nil {
		fmt.Printf("ERROR - revenues.go/interface/InsertIntoRevenues -  %s\n", err)
		return 0
	}
	fmt.Printf("Success! Execution time: %d milliseconds\n", time.Now().UnixMilli()-execTime)
	return output
}
