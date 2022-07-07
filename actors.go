package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//
// ---------------------------------------- CONFIG SECTION ----------------------------------------
// ---------------------------------------- CONFIG SECTION ----------------------------------------
// ---------------------------------------- CONFIG SECTION ----------------------------------------
//
// Defining a json-friendly struct for rows from 'actors' table
type Actors struct {
	ActorID     *int    `json:"actorID"`
	FirstName   *string `json:"firstName"`
	LastName    *string `json:"lastName"`
	Gender      *string `json:"gender"`
	DateOfBirth *string `json:"dateOfBirth"`
}

const (
	actorStatementSingle        = "SELECT * FROM actors WHERE actor_id = $1"
	actorStatementAll           = "SELECT * FROM actors"
	actorStatementInsertInto    = "INSERT INTO actors (first_name, last_name, gender, date_of_birth) VALUES ($1, $2, $3, $4) RETURNING actor_id"
	actorStatementDeleteFrom    = "DELETE FROM actors WHERE actor_id = $1 RETURNING actor_id"
	actorStatementCheckIfExists = "SELECT count(actor_id) FROM actors WHERE actor_id = $1"
)

//
// ---------------------------------------- SERVICE SECTION ----------------------------------------
// ---------------------------------------- SERVICE SECTION ----------------------------------------
// ---------------------------------------- SERVICE SECTION ----------------------------------------
//
// Get a single actor from the database, encode it into a json
// object and send it as a response. Log the activity into console
func getActor(w http.ResponseWriter, r *http.Request) {
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
	if exists, err := checkIfExists(id, actorStatementCheckIfExists); err != nil || !exists {
		w.WriteHeader(statusCodeNotFound.status)
		w.Write([]byte(statusCodeNotFound.message + " - item does not exist!"))
		fmt.Printf("ERROR: actors.go/checkIfExists:   %s\n", err)
		return
	}
	// CHECK if there was an error during query
	actorRow, err := selectItem(id, actorStatementSingle)
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
	var actors []Actors
	// CHECK if there was an error during query
	actorRows, err := selectAllItems(actorStatementAll)
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

// Find a single actor in a database and delete it permanently
func deleteActor(w http.ResponseWriter, r *http.Request) {
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
	if exists, err := checkIfExists(id, actorStatementCheckIfExists); err != nil || !exists {
		w.WriteHeader(statusCodeNotFound.status)
		w.Write([]byte(statusCodeNotFound.message + " - item does not exist!"))
		fmt.Printf("ERROR: actors.go/checkIfExists:   %s\n", err)
		return
	}

	// Sending DELETE FROM query
	deletedID, err := deleteItem(id, actorStatementDeleteFrom)
	if err != nil {
		fmt.Printf("ERROR: actors.go/deleteItem: %s\n", err)
		return
	}
	w.WriteHeader(statusCodeItemDeleted.status)
	w.Write([]byte(fmt.Sprintf("%s - %d", statusCodeItemDeleted.message, deletedID)))
	w.Header().Set("Content-Type", "application/json")

	fmt.Printf("DELETE FROM | actors | %d\n", deletedID)
}

//
// ---------------------------------------- INTERFACE SECTION ----------------------------------------
// ---------------------------------------- INTERFACE SECTION ----------------------------------------
// ---------------------------------------- INTERFACE SECTION ----------------------------------------
//
// rowScanner method for scanning an 'Actors' object from a single row
func (inputRow *SqlRow) ScanActor() Actors {
	var output Actors
	err := inputRow.row.Scan(&output.ActorID, &output.FirstName, &output.LastName, &output.Gender, &output.DateOfBirth)
	if err != nil {
		fmt.Printf("ERROR - actors.go/interface -  %s\n", err)
		return output
	}
	return output
}

// rowScanner method for scanning an 'Actors' object from multiple rows
func (inputRow *SqlRows) ScanAllActors() Actors {
	var output Actors
	err := inputRow.rows.Scan(&output.ActorID, &output.FirstName, &output.LastName, &output.Gender, &output.DateOfBirth)
	if err != nil {
		fmt.Printf("ERROR - actors.go/interface -  %s\n", err)
		return output
	}
	return output
}
