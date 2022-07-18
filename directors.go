package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

//
// ---------------------------------------- CONFIG SECTION ----------------------------------------
// ---------------------------------------- CONFIG SECTION ----------------------------------------
// ---------------------------------------- CONFIG SECTION ----------------------------------------
//
// Defining a json-friendly struct for rows from 'directors' table
type Directors struct {
	DirectorID  *int    `json:"directorID"`
	FirstName   *string `json:"firstName"`
	LastName    *string `json:"lastName"`
	DateOfBirth *string `json:"dateOfBirth"`
	Nationality *string `json:"nationality"`
}

const (
	directorStatementSingle     = "SELECT * FROM directors WHERE director_id = $1"
	directorStatementAll        = "SELECT * FROM directors"
	directorStatementInsertInto = "INSERT INTO directors (first_name, last_name, date_of_birth, nationality) VALUES ($1, $2, $3, $4) RETURNING director_id"
	directorStatementDeleteFrom = "DELETE FROM directors WHERE director_id = $1 RETURNING director_id"
	directorStatementCheckForID = "SELECT count(director_id) FROM directors WHERE director_id = $1"
)

//
// ---------------------------------------- SERVICE SECTION ----------------------------------------
// ---------------------------------------- SERVICE SECTION ----------------------------------------
// ---------------------------------------- SERVICE SECTION ----------------------------------------
//
// Get a single director from the database, encode it into a json
// object and send it as a response. Log the activity into console
func getDirector(w http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	// CHECK for a valid ID parameter
	if err != nil {
		w.WriteHeader(statusCodeBadRequest.status)
		w.Write([]byte(statusCodeBadRequest.message + " - ID not an integer!"))
		fmt.Println("ERROR: directors.go/strconv.Atoi():   Not integer!")
		return
	}
	// CHECK if item of ID exists in the database
	if exists, err := checkIfExists(id, directorStatementCheckForID); err != nil || !exists {
		w.WriteHeader(statusCodeNotFound.status)
		w.Write([]byte(statusCodeNotFound.message + " - item does not exist!"))
		fmt.Printf("ERROR: directors.go/checkIfExists:   %s\n", err)
		return
	}
	// CHECK if there was an error during query
	directorRow, err := selectItem(id, directorStatementSingle)
	if err != nil {
		w.WriteHeader(statusCodeInternalError.status)
		w.Write([]byte(statusCodeInternalError.message + " - could not fetch data!"))
		fmt.Printf("ERROR: directors.go/selectItem:   %s\n", err)
		return
	}
	// SCAN row data into json-able object
	director := directorRow.ScanDirector()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(director)

	fmt.Printf("SELECT | directors | %d\n", *director.DirectorID)
}

// Get all rows of directors from the database, encode it into a json
// object and send it as a response. Log the activity into console
func getAllDirectors(w http.ResponseWriter, r *http.Request) {
	var directors []Directors
	// CHECK if there was an error during query
	directorRows, err := selectAllItems(directorStatementAll)
	if err != nil {
		w.WriteHeader(statusCodeInternalError.status)
		w.Write([]byte(statusCodeInternalError.message + " - could not fetch data!"))
		fmt.Printf("ERROR: directors.go/selectItem:   %s\n", err)
		return
	}
	// SCAN rows data into json-able object
	for directorRows.rows.Next() {
		director := directorRows.ScanAllDirectors()
		directors = append(directors, director)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(directors)

	fmt.Println("SELECT | directors | ALL")
}

// Create a single director in a database based on input JSON object
func postDirector(w http.ResponseWriter, r *http.Request) {
	var director Directors
	var statement SqlStmt
	var err error
	json.NewDecoder(r.Body).Decode(&director)
	// Add error handling
	statement.statement, err = insertItem(directorStatementInsertInto)
	if err != nil {
		w.WriteHeader(statusCodeQueryError.status)
		w.Write([]byte(statusCodeQueryError.message))
		fmt.Printf("ERROR: directors.go/insertItem:   %s\n", err)
		return
	}
	id := statement.InsertIntoDirectors(director)
	if id <= 0 {
		w.WriteHeader(statusCodeQueryError.status)
		w.Write([]byte(statusCodeQueryError.message + " - Could not execute query!"))
		fmt.Printf("ERROR: directors.go/InsertIntoDirectors:   %s\n", err)
		return
	}
	statement.statement.Close()
	director.DirectorID = &id

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(director)

	fmt.Printf("INSERT INTO | directors | %d\n", *director.DirectorID)
}

// Find a single director in a database and delete it permanently
func deleteDirector(w http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	// CHECK for a valid ID parameter
	if err != nil {
		w.WriteHeader(statusCodeBadRequest.status)
		w.Write([]byte(statusCodeBadRequest.message + " - ID not an integer!"))
		fmt.Println("ERROR: directors.go/strconv.Atoi():   Not integer!")
		return
	}
	// CHECK if item of ID exists in the database
	if exists, err := checkIfExists(id, directorStatementCheckForID); err != nil || !exists {
		w.WriteHeader(statusCodeNotFound.status)
		w.Write([]byte(statusCodeNotFound.message + " - item does not exist!"))
		fmt.Printf("ERROR: directors.go/checkIfExists:   %s\n", err)
		return
	}

	// Sending DELETE FROM query
	deletedID, err := deleteItem(id, directorStatementDeleteFrom)
	if err != nil {
		fmt.Printf("ERROR: directors.go/deleteItem: %s\n", err)
		return
	}
	w.WriteHeader(statusCodeItemDeleted.status)
	w.Write([]byte(fmt.Sprintf("%s - %d", statusCodeItemDeleted.message, deletedID)))
	w.Header().Set("Content-Type", "application/json")

	fmt.Printf("DELETE FROM | directors | %d\n", deletedID)
}

//
// ---------------------------------------- INTERFACE SECTION ----------------------------------------
// ---------------------------------------- INTERFACE SECTION ----------------------------------------
// ---------------------------------------- INTERFACE SECTION ----------------------------------------
//
// rowScanner method for scanning a 'Directors' object
func (inputRow *SqlRow) ScanDirector() Directors {
	var output Directors
	err := inputRow.row.Scan(&output.DirectorID, &output.FirstName, &output.LastName, &output.DateOfBirth, &output.Nationality)
	if err != nil {
		fmt.Printf("ERROR - directors.go/interface -  %s\n", err)
		return output
	}
	return output
}

// rowScanner method for scanning an 'Directors' object from multiple rows
func (inputRow *SqlRows) ScanAllDirectors() Directors {
	var output Directors
	err := inputRow.rows.Scan(&output.DirectorID, &output.FirstName, &output.LastName, &output.DateOfBirth, &output.Nationality)
	if err != nil {
		fmt.Printf("ERROR - directors.go/interface -  %s\n", err)
		return output
	}
	return output
}

// rowScanner method for inserting an 'Directors' object into the database as a single row
func (queryStatement *SqlStmt) InsertIntoDirectors(inputRow Directors) int {
	var output int
	execTime := time.Now().UnixMilli()
	err := queryStatement.statement.QueryRow(inputRow.FirstName, inputRow.LastName, inputRow.DateOfBirth, inputRow.Nationality).Scan(&output)
	if err != nil {
		fmt.Printf("ERROR - directors.go/interface -  %s\n", err)
		return 0
	}
	fmt.Printf("Success! Execution time: %d milliseconds\n", time.Now().UnixMilli()-execTime)
	return output
}
