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
// Defining a json-friendly struct for rows from 'movies' table
type Movies struct {
	MovieID        *int    `json:"movieID"`
	MovieName      *string `json:"movieName"`
	MovieLength    *string `json:"movieLength"`
	MovieLang      *string `json:"movieLang"`
	ReleaseDate    *string `json:"releaseDate"`
	AgeCertificate *string `json:"ageCertificate"`
	DirectorID     *int    `json:"directorID"`
}

const (
	movieStatementSingle        = "SELECT * FROM movies WHERE movie_id = $1"
	movieStatementAll           = "SELECT * FROM movies"
	movieStatementInsertInto    = "INSERT INTO movies (movie_name, movie_length, movie_lang, release_date, age_certificate, director_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING movie_id"
	movieStatementDeleteFrom    = "DELETE FROM movies WHERE movie_id = $1 RETURNING movie_id"
	movieStatementCheckIfExists = "SELECT count(movie_id) FROM movies WHERE movie_id = $1"
)

//
// -------------------- SERVICE SECTION --------------------
// -------------------- SERVICE SECTION --------------------
// -------------------- SERVICE SECTION --------------------
//
// Get a single movie from the database, encode it into a json
// object and send it as a response. Log the activity into console
func getMovie(w http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	// CHECK for a valid ID parameter
	if err != nil {
		w.WriteHeader(statusCodeBadRequest.status)
		w.Write([]byte(statusCodeBadRequest.message + " - ID not an integer!"))
		fmt.Println("ERROR: movies.go/strconv.Atoi():   Not integer!")
		return
	}
	// CHECK if item of ID exists in the database
	if exists, err := checkIfExists(id, movieStatementCheckIfExists); err != nil || !exists {
		w.WriteHeader(statusCodeNotFound.status)
		w.Write([]byte(statusCodeNotFound.message + " - item does not exist!"))
		fmt.Printf("ERROR: movies.go/checkIfExists:   %s\n", err)
		return
	}
	// CHECK if there was an error during query
	movieRow, err := selectItem(id, movieStatementSingle)
	if err != nil {
		w.WriteHeader(statusCodeInternalError.status)
		w.Write([]byte(statusCodeInternalError.message + " - could not fetch data!"))
		fmt.Printf("ERROR: movies.go/selectItem:   %s\n", err)
		return
	}
	// SCAN row data into json-able object
	movie := movieRow.ScanMovie()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movie)

	fmt.Printf("SELECT | movies | %d\n", *movie.MovieID)
}
