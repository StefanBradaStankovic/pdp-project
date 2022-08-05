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

// Defining a json-friendly struct for rows from 'actors' table
type Actors struct {
	ActorID     *int    `json:"actorID"`
	FirstName   *string `json:"firstName"`
	LastName    *string `json:"lastName"`
	Gender      *string `json:"gender"`
	DateOfBirth *string `json:"dateOfBirth"`
}

var actorStatement = queryStatement{
	"SELECT actor_id, first_name, last_name, gender, date_of_birth FROM actors WHERE actor_id = $1 AND is_visible", //selectSingle
	"SELECT actor_id, first_name, last_name, gender, date_of_birth FROM actors WHERE is_visible",                   //selectAll
	"INSERT INTO actors (first_name, last_name, gender, date_of_birth) VALUES ($1, $2, $3, $4) RETURNING actor_id", //insertInto
	"DELETE FROM actors WHERE actor_id = $1 RETURNING actor_id",                                                    //deleteFrom
	"UPDATE actors SET is_visible = false WHERE actor_id = $1",                                                     //updateVisible
	"SELECT count(actor_id) FROM actors WHERE actor_id = $1 AND is_visible"}                                        //checkForID

// ---------------------------------------- SERVICE SECTION ----------------------------------------

// Get a single actor from the database, encode it into a json
// object and send it as a response. Log the activity into console
func getActor(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)
	var idParam string = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	// CHECK for a valid ID parameter
	if err != nil {
		w.WriteHeader(statusCodeBadRequest.status)
		w.Write([]byte(statusCodeBadRequest.message + " - ID not an integer!"))
		fmt.Println("ERROR: actors.go/strconv.Atoi():   Not integer!")
		return
	}
	// CHECK if item of ID exists in the database
	if exists, err := checkIfExists(id, actorStatement.checkForID); err != nil || !exists {
		w.WriteHeader(statusCodeNotFound.status)
		w.Write([]byte(statusCodeNotFound.message + " - item does not exist!"))
		return
	}
	// CHECK if there was an error during query
	actorRow, err := selectItem(id, actorStatement.selectSingle)
	if err != nil {
		w.WriteHeader(statusCodeInternalError.status)
		w.Write([]byte(statusCodeInternalError.message + " - could not fetch data!"))
		fmt.Printf("ERROR: actors.go/selectItem:   %s\n", err)
		return
	}
	// SCAN row data into json-able object
	actor := actorRow.ScanActor()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actor)

	fmt.Printf("SELECT | actors | %d\n", *actor.ActorID)
}

// Get all rows of actors from the database, encode it into a json
// object and send it as a response. Log the activity into console
func getAllActors(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)
	var actors []Actors
	// CHECK if there was an error during query
	actorRows, err := selectAllItems(actorStatement.selectAll)
	if err != nil {
		w.WriteHeader(statusCodeInternalError.status)
		w.Write([]byte(statusCodeInternalError.message + " - could not fetch data!"))
		fmt.Printf("ERROR: actors.go/selectItem:   %s\n", err)
		return
	}

	// SCAN rows data into json-able object
	for actorRows.rows.Next() {
		actor := actorRows.ScanAllActors()
		actors = append(actors, actor)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actors)

	fmt.Println("SELECT | actors | ALL")
}

// Create a single actor in a database based on input JSON object
func postActor(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)
	var actor Actors
	var statement SqlStmt
	var err error
	json.NewDecoder(r.Body).Decode(&actor)
	// Add error handling
	statement.statement, err = insertItem(actorStatement.insertInto)
	if err != nil {
		w.WriteHeader(statusCodeQueryError.status)
		w.Write([]byte(statusCodeQueryError.message))
		fmt.Printf("ERROR: actors.go/insertItem:   %s\n", err)
		return
	}
	id := statement.InsertIntoActors(actor)
	if id <= 0 {
		w.WriteHeader(statusCodeQueryError.status)
		w.Write([]byte(statusCodeQueryError.message + " - Could not execute query!"))
		return
	}
	statement.statement.Close()
	actor.ActorID = &id

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actor)

	fmt.Printf("INSERT INTO | actors | %d\n", *actor.ActorID)
}

// Find a single actor in a database and delete it permanently
func deleteActor(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)
	var idParam string = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	// CHECK for a valid ID parameter
	if err != nil {
		w.WriteHeader(statusCodeBadRequest.status)
		w.Write([]byte(statusCodeBadRequest.message + " - ID not an integer!"))
		fmt.Println("ERROR: actors.go/strconv.Atoi():   Not integer!")
		return
	}
	// CHECK if item of ID exists in the database
	if exists, err := checkIfExists(id, actorStatement.checkForID); err != nil || !exists {
		w.WriteHeader(statusCodeNotFound.status)
		w.Write([]byte(statusCodeNotFound.message + " - item does not exist!"))
		return
	}

	// Sending DELETE FROM query
	err = deleteItem(id, actorStatement.updateVisible)
	if err != nil {
		fmt.Printf("ERROR: actors.go/deleteItem: %s\n", err)
		return
	}
	w.WriteHeader(statusCodeItemDeleted.status)
	w.Write([]byte(fmt.Sprintf("%s - %d", statusCodeItemDeleted.message, id)))
	w.Header().Set("Content-Type", "application/json")

	fmt.Printf("DELETE FROM | actors | %d\n", id)
}

// ---------------------------------------- INTERFACE SECTION ----------------------------------------

// rowScanner method for scanning an 'Actors' object from a single row
func (inputRow *SqlRow) ScanActor() Actors {
	var output Actors
	err := inputRow.row.Scan(&output.ActorID, &output.FirstName, &output.LastName, &output.Gender, &output.DateOfBirth)
	if err != nil {
		fmt.Printf("ERROR - actors.go/interface/ScanActors -  %s\n", err)
		return output
	}
	return output
}

// rowScanner method for scanning an 'Actors' object from multiple rows
func (inputRow *SqlRows) ScanAllActors() Actors {
	var output Actors
	err := inputRow.rows.Scan(&output.ActorID, &output.FirstName, &output.LastName, &output.Gender, &output.DateOfBirth)
	if err != nil {
		fmt.Printf("ERROR - actors.go/interface/ScanAllActors -  %s\n", err)
		return output
	}
	return output
}

// rowScanner method for inserting an 'Actors' object into the database as a single row
func (queryStatement *SqlStmt) InsertIntoActors(inputRow Actors) int {
	var output int
	execTime := time.Now().UnixMilli()
	err := queryStatement.statement.QueryRow(inputRow.FirstName, inputRow.LastName, inputRow.Gender, inputRow.DateOfBirth).Scan(&output)
	if err != nil {
		fmt.Printf("ERROR - actors.go/interface/InsertIntoActors -  %s\n", err)
		return 0
	}
	fmt.Printf("Success! Execution time: %d milliseconds\n", time.Now().UnixMilli()-execTime)
	return output
}
