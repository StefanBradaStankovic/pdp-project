package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	// INIT START
	fmt.Printf("Starting up the server . . .\n")
	router := mux.NewRouter()
	psqlConnect = setDBConnection()
	db = dbConnect()
	defer db.Close()
	// INIT END

	router.HandleFunc("/actors/{id}", getActor).Methods("GET")       // GET a single actor by ID
	router.HandleFunc("/directors/{id}", getDirector).Methods("GET") // GET a single director by ID
	router.HandleFunc("/movies/{id}", getMovie).Methods("GET")       // GET a single movie by ID
	router.HandleFunc("/revenues/{id}", getRevenues).Methods("GET")  // GET a single set of revenues by ID

	router.HandleFunc("/actors", getAllActors).Methods("GET")       // GET all actors
	router.HandleFunc("/directors", getAllDirectors).Methods("GET") // GET all directors
	router.HandleFunc("/movies", getAllMovies).Methods("GET")       // GET all movies
	router.HandleFunc("/revenues", getAllRevenues).Methods("GET")   // GET all movie revenues

	router.HandleFunc("/actors", postActor).Methods("POST")       // POST a single actor
	router.HandleFunc("/directors", postDirector).Methods("POST") // POST a single director
	router.HandleFunc("/movies", postMovie).Methods("POST")       // POST a single movie
	router.HandleFunc("/revenues", postRevenues).Methods("POST")  // POST a single set of revenues

	router.HandleFunc("/actors/{id}", deleteActor).Methods("DELETE")       // DELETE a single actor by ID
	router.HandleFunc("/directors/{id}", deleteDirector).Methods("DELETE") // DELETE a single director by ID
	router.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")       // DELETE a single movie by ID
	router.HandleFunc("/revenues/{id}", deleteRevenues).Methods("DELETE")  // DELETE a single set of revenues by ID

	fmt.Printf("Server up and running!\n")
	http.ListenAndServe(":5000", router)
}
