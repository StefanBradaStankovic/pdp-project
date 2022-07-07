package main

/*

//
// ---------- GET SECTION ----------
//
// Query one row by ID from table 'actors'

// Query one row by ID from table 'directors'
func getDirector(w http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("ID could not be converted to integer"))
		fmt.Println("ERROR: bad ID: Not integer!")
		return
	}

	if id <= 0 || id > TablesLimit.DirectorsLimit {
		w.WriteHeader(400)
		w.Write([]byte("ID out of range"))
		fmt.Println("ERROR: bad ID: Out of range!")
		return
	}

	fmt.Printf("Requested Director of ID %d\n", id)

	director, err := selectDirector(id)
	if err != nil {
		fmt.Printf("ERROR: bad query: %s", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(NillableToDirectors(director))

	fmt.Printf("Returning:\n%v\n", NillableToDirectors(director))
}

// Query one row by ID from table 'movies'
func getMovie(w http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("ID could not be converted to integer"))
		fmt.Println("ERROR: bad ID: Not integer!")
		return
	}

	if id <= 0 || id > TablesLimit.MoviesLimit {
		w.WriteHeader(400)
		w.Write([]byte("ID out of range"))
		fmt.Println("ERROR: bad ID: Out of range!")
		return
	}

	fmt.Printf("Requested Movie of ID %d\n", id)

	movie, err := selectMovie(id)
	if err != nil {
		fmt.Printf("ERROR: bad query: %s", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(NillableToMovies(movie))

	fmt.Printf("Returning:\n%v\n", NillableToMovies(movie))
}

// Query one row by ID from table 'directors'
func getRevenues(w http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("ID could not be converted to integer"))
		fmt.Println("ERROR: bad ID: Not integer!")
		return
	}

	if id <= 0 || id > TablesLimit.MovieRevenuesLimit {
		w.WriteHeader(400)
		w.Write([]byte("ID out of range"))
		fmt.Println("ERROR: bad ID: Out of range!")
		return
	}

	fmt.Printf("Requested Revenue of ID %d\n", id)

	revenue, err := selectMovieRevenues(id)
	if err != nil {
		fmt.Printf("ERROR: bad query: %s\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(NillableToMovieRevenues(revenue))

	fmt.Printf("Returning:\n%v\n", NillableToMovieRevenues(revenue))
}

// Query all rows from table 'actors'
func getAllActors(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Requested Actor ALL ROWS")

	actors, err := selectAllActors()
	if err != nil {
		fmt.Printf("ERROR: %e\n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actors)
}

// Query all rows from table 'directors'
func getAllDirectors(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Requested Directors ALL ROWS")

	directors, err := selectAllDirectors()
	if err != nil {
		fmt.Printf("ERROR: %e\n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(directors)
}

// Query all rows from table 'movies'
func getAllMovies(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Requested Movies ALL ROWS")

	movies, err := selectAllMovies()
	if err != nil {
		fmt.Printf("ERROR: %e\n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

// Query all rows from table 'movie_revenues'
func getAllMovieRevenues(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Requested Revenues ALL ROWS")

	movieRevenues, err := selectAllMovieRevenues()
	if err != nil {
		fmt.Printf("ERROR: %e\n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movieRevenues)
}

//
// ---------- POST SECTION ----------
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
}

//
// ---------- DELETE SECTION ----------
//
// Delete one row by ID from table 'actors'
func deleteActor(w http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("ID could not be converted to integer"))
		fmt.Println("ERROR: bad ID: Not integer!")
		return
	}

	// Check if ID exists instead of current method
	if id <= 0 || id > TablesLimit.ActorsLimit {
		w.WriteHeader(400)
		w.Write([]byte("ID out of range"))
		fmt.Println("ERROR: bad ID: Out of range!")
		return
	}

	fmt.Printf("Deleting actor of ID %d\n", id)

	deletedId, err := deleteFromActors(id)
	if err != nil {
		fmt.Printf("ERROR: bad query: %s\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deletedId)

	fmt.Printf("Deleting:\n%d\n", deletedId)

}

// Delete one row by ID from table 'directors'
func deleteDirector(w http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("ID could not be converted to integer"))
		fmt.Println("ERROR: bad ID: Not integer!")
		return
	}

	if id <= 0 || id > TablesLimit.DirectorsLimit {
		w.WriteHeader(400)
		w.Write([]byte("ID out of range"))
		fmt.Println("ERROR: bad ID: Out of range!")
		return
	}

	fmt.Printf("Deleting director of ID %d\n", id)

	deletedId, err := deleteFromDirectors(id)
	if err != nil {
		fmt.Printf("ERROR: bad query: %s\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deletedId)

	fmt.Printf("Deleting:\n%d\n", deletedId)

}

// Delete one row by ID from table 'movies'
func deleteMovie(w http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("ID could not be converted to integer"))
		fmt.Println("ERROR: bad ID: Not integer!")
		return
	}

	if id <= 0 || id > TablesLimit.MoviesLimit {
		w.WriteHeader(400)
		w.Write([]byte("ID out of range"))
		fmt.Println("ERROR: bad ID: Out of range!")
		return
	}

	fmt.Printf("Deleting movie of ID %d\n", id)

	deletedId, err := deleteFromMovies(id)
	if err != nil {
		fmt.Printf("ERROR: bad query: %s\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deletedId)

	fmt.Printf("Deleting:\n%d\n", deletedId)

}

// Delete one row by ID from table 'actors'
func deleteMovieRevenues(w http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("ID could not be converted to integer"))
		fmt.Println("ERROR: bad ID: Not integer!")
		return
	}

	if id <= 0 || id > TablesLimit.MovieRevenuesLimit {
		w.WriteHeader(400)
		w.Write([]byte("ID out of range"))
		fmt.Println("ERROR: bad ID: Out of range!")
		return
	}

	fmt.Printf("Deleting actor of ID %d\n", id)

	deletedId, err := deleteFromMovieRevenues(id)
	if err != nil {
		fmt.Printf("ERROR: bad query: %s\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deletedId)

	fmt.Printf("Deleting:\n%d\n", deletedId)

}

*/
