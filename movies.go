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

var movieStatement = queryStatement{
	"SELECT movie_id, movie_name, movie_length, movie_lang, release_date, age_certificate, director_id FROM movies WHERE movie_id = $1 AND is_visible",         //selectSingle
	"SELECT movie_id, movie_name, movie_length, movie_lang, release_date, age_certificate, director_id FROM movies WHERE is_visible",                           //selectAll
	"INSERT INTO movies (movie_name, movie_length, movie_lang, release_date, age_certificate, director_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING movie_id", //insertInto
	"DELETE FROM movies WHERE movie_id = $1 RETURNING movie_id",                                                                                                //deleteFrom
	"UPDATE movies SET is_visible = false WHERE movie_id = $1",                                                                                                 //updateVisible
	"SELECT count(movie_id) FROM movies WHERE movie_id = $1 AND is_visible"}                                                                                    //checkForID

// ---------------------------------------- SERVICE SECTION ----------------------------------------

// Get a single movie from the database, encode it into a json
// object and send it as a response. Log the activity into console
func getMovie(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)
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
	if exists, err := checkIfExists(id, movieStatement.checkForID); err != nil || !exists {
		w.WriteHeader(statusCodeNotFound.status)
		w.Write([]byte(statusCodeNotFound.message + " - item does not exist!"))
		return
	}
	// CHECK if there was an error during query
	movieRow, err := selectItem(id, movieStatement.selectSingle)
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

// Get all rows of movies from the database, encode it into a json
// object and send it as a response. Log the activity into console
func getAllMovies(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)
	var movies []Movies
	// CHECK if there was an error during query
	moviesRows, err := selectAllItems(movieStatement.selectAll)
	if err != nil {
		w.WriteHeader(statusCodeInternalError.status)
		w.Write([]byte(statusCodeInternalError.message + " - could not fetch data!"))
		fmt.Printf("ERROR: movies.go/selectItem:   %s\n", err)
		return
	}
	// SCAN rows data into json-able object
	for moviesRows.rows.Next() {
		movie := moviesRows.ScanAllMovies()
		movies = append(movies, movie)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)

	fmt.Println("SELECT | movies | ALL")
}

// Create a single movie in a database based on input JSON object
func postMovie(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)
	var movie Movies
	var statement SqlStmt
	var err error
	json.NewDecoder(r.Body).Decode(&movie)
	// Add error handling
	statement.statement, err = insertItem(movieStatement.insertInto)
	if err != nil {
		w.WriteHeader(statusCodeQueryError.status)
		w.Write([]byte(statusCodeQueryError.message))
		fmt.Printf("ERROR: movies.go/insertItem:   %s\n", err)
		return
	}
	id := statement.InsertIntoMovies(movie)
	if id <= 0 {
		w.WriteHeader(statusCodeQueryError.status)
		w.Write([]byte(statusCodeQueryError.message + " - Could not execute query!"))
		fmt.Printf("ERROR: movies.go/InsertIntoMovies:   %s\n", err)
		return
	}
	statement.statement.Close()
	movie.MovieID = &id

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movie)

	fmt.Printf("INSERT INTO | movies | %d\n", *movie.MovieID)
}

// Find a single movie in a database and delete it permanently
func deleteMovie(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)
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
	if exists, err := checkIfExists(id, movieStatement.checkForID); err != nil || !exists {
		w.WriteHeader(statusCodeNotFound.status)
		w.Write([]byte(statusCodeNotFound.message + " - item does not exist!"))
		return
	}

	// Sending DELETE FROM query
	err = deleteItem(id, movieStatement.deleteFrom)
	if err != nil {
		fmt.Printf("ERROR: movies.go/deleteItem: %s\n", err)
		return
	}
	w.WriteHeader(statusCodeItemDeleted.status)
	w.Write([]byte(fmt.Sprintf("%s - %d", statusCodeItemDeleted.message, id)))
	w.Header().Set("Content-Type", "application/json")

	fmt.Printf("DELETE FROM | movies | %d\n", id)
}

// ---------------------------------------- INTERFACE SECTION ----------------------------------------

// rowScanner method for scanning a 'Movies' object
func (inputRow *SqlRow) ScanMovie() Movies {
	var output Movies
	err := inputRow.row.Scan(&output.MovieID, &output.MovieName, &output.MovieLength, &output.MovieLang, &output.ReleaseDate, &output.AgeCertificate, &output.DirectorID)
	if err != nil {
		fmt.Printf("ERROR - movies.go/interface/ScanMovie -  %s\n", err)
		return output
	}
	return output
}

// rowScanner method for scanning an 'Movies' object from multiple rows
func (inputRow *SqlRows) ScanAllMovies() Movies {
	var output Movies
	err := inputRow.rows.Scan(&output.MovieID, &output.MovieName, &output.MovieLength, &output.MovieLang, &output.ReleaseDate, &output.AgeCertificate, &output.DirectorID)
	if err != nil {
		fmt.Printf("ERROR - movies.go/interface/ScanAllMovies -  %s\n", err)
		return output
	}
	return output
}

// rowScanner method for inserting an 'Movies' object into the database as a single row
func (queryStatement *SqlStmt) InsertIntoMovies(inputRow Movies) int {
	var output int
	execTime := time.Now().UnixMilli()
	err := queryStatement.statement.QueryRow(inputRow.MovieName, inputRow.MovieLength, inputRow.MovieLang, inputRow.ReleaseDate, inputRow.AgeCertificate, inputRow.DirectorID).Scan(&output)
	if err != nil {
		fmt.Printf("ERROR - actors.go/interface/InsertIntoMovies -  %s\n", err)
		return 0
	}
	fmt.Printf("Success! Execution time: %d milliseconds\n", time.Now().UnixMilli()-execTime)
	return output
}
