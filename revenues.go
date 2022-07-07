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
// Defining a json-friendly struct for rows from 'movie_revenues' table
type MovieRevenues struct {
	RevenueID            *int     `json:"revenueID"`
	MovieID              *int     `json:"movieID"`
	DomesticTakings      *float64 `json:"domesticTakings"`
	InternationalTakings *float64 `json:"internationalTakings"`
}

const (
	revenuesStatementSingle        = "SELECT * FROM movie_revenues WHERE revenue_id = $1"
	revenuesStatementAll           = "SELECT * FROM movie_revenues"
	revenuesStatementInsertInto    = "INSERT INTO movie_revenues (revenue_id, movie_id, domestic_takings, international_takings) VALUES ($1, $2, $3, $4) RETURNING revenue_id"
	revenuesStatementDeleteFrom    = "DELETE FROM movie_revenues WHERE revenue_id = $1 RETURNING revenue_id"
	revenuesStatementCheckIfExists = "SELECT count(revenue_id) FROM movie_revenues WHERE revenue_id = $1"
)

//
// -------------------- SERVICE SECTION --------------------
// -------------------- SERVICE SECTION --------------------
// -------------------- SERVICE SECTION --------------------
//
// Get a single set of revenues from the database, encode it into a json
// object and send it as a response. Log the activity into console
func getRevenues(w http.ResponseWriter, r *http.Request) {
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
	if exists, err := checkIfExists(id, revenuesStatementCheckIfExists); err != nil || !exists {
		w.WriteHeader(statusCodeNotFound.status)
		w.Write([]byte(statusCodeNotFound.message + " - item does not exist!"))
		fmt.Printf("ERROR: revenues.go/checkIfExists:   %s\n", err)
		return
	}
	// CHECK if there was an error during query
	revenueRow, err := selectItem(id, revenuesStatementSingle)
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

	fmt.Printf("SELECT | revenues | %d\n", *revenue.RevenueID)
}
