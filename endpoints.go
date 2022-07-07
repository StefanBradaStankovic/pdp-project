package main

/*//
// ---------------------------------------- POST SECTION ----------------------------------------
// ---------------------------------------- POST SECTION ----------------------------------------
// ---------------------------------------- POST SECTION ----------------------------------------
//
// Create a new entry in 'actors' table
func postActor(w http.ResponseWriter, r *http.Request) {
	var actor Actors
	json.NewDecoder(r.Body).Decode(&actor)
	fmt.Printf("Received object: %v\n", actor)

	// Add error handling
	actor.ActorID = insertIntoActors(actor)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actor)
}

// Create a new entry in 'directors' table
func postDirector(w http.ResponseWriter, r *http.Request) {
	var director Directors
	json.NewDecoder(r.Body).Decode(&director)
	fmt.Printf("Received object: %v\n", director)

	director.DirectorID = insertIntoDirectors(director)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(director)
}

// Create a new entry in 'movies' table
func postMovie(w http.ResponseWriter, r *http.Request) {
	var movie Movies
	json.NewDecoder(r.Body).Decode(&movie)
	fmt.Printf("Received object: %v\n", movie)

	movie.MovieID = insertIntoMovies(movie)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movie)
}

// Create a new entry in 'movie_revenues' table
func postMovieRevenues(w http.ResponseWriter, r *http.Request) {
	var revenues MovieRevenues
	json.NewDecoder(r.Body).Decode(&revenues)
	fmt.Printf("Received object: %v\n", revenues)

	revenues.RevenueID = insertIntoMovieRevenues(revenues)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(revenues)
}*/
