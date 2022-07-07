package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//
// -------------------- CONFIG SECTION --------------------
// -------------------- CONFIG SECTION --------------------
// -------------------- CONFIG SECTION --------------------
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
// -------------------- SERVICE SECTION --------------------
// -------------------- SERVICE SECTION --------------------
// -------------------- SERVICE SECTION --------------------
//
// Get a single actor from the database, encode it into a json
// object and send it as a response. Log the activity into console
func getActor(w http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil { // CHECK for a valid ID parameter
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
