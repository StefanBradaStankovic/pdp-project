package main

/*
//
// ---------- SELECT SECTION ----------
//
// Select a single row from 'TABLE_NAME' using 'ITEM_ID'
func selectItem(id int, queryString string) (SqlRow, error) {
	var item SqlRow
	statement, err := db.Prepare(queryString)
	if err != nil {
		return item, err
	}

	item.row = statement.QueryRow(id)
	err = item.row.Err()
	if err != nil {
		return item, err
	}

	return item, err
}



// Sellect a single row from 'directors' table using director_id
func selectDirector(id int) (NillableDirectors, error) {
	var director NillableDirectors
	statement, err := db.Prepare("SELECT * FROM directors WHERE director_id = $1")
	if err != nil {
		return director, err
	}

	execTime := time.Now().UnixMilli()
	err = statement.QueryRow(id).Scan(&director.DirectorID, &director.FirstName, &director.LastName, &director.Nationality, &director.DateOfBirth)
	if err != nil {
		return director, err
	}
	fmt.Printf("Success! Execution time: %d milliseconds\n", time.Now().UnixMilli()-execTime)

	return director, err
}

// Sellect a single row from 'movies' table using movie_id
func selectMovie(id int) (NillableMovies, error) {
	var movie NillableMovies
	statement, err := db.Prepare("SELECT * FROM movies WHERE movie_id = $1")
	if err != nil {
		return movie, err
	}

	execTime := time.Now().UnixMilli()
	err = statement.QueryRow(id).Scan(&movie.MovieID, &movie.MovieName, &movie.MovieLength, &movie.MovieLang, &movie.ReleaseDate, &movie.AgeCertificate, &movie.DirectorID)
	if err != nil {
		return movie, err
	}
	fmt.Printf("Success! Execution time: %d milliseconds\n", time.Now().UnixMilli()-execTime)

	return movie, err
}

// Sellect a single row from 'movie_revenues' table using revenue_id
func selectMovieRevenues(id int) (NillableMovieRevenues, error) {
	var revenues NillableMovieRevenues
	statement, err := db.Prepare("SELECT * FROM movie_revenues WHERE revenue_id = $1")
	if err != nil {
		return revenues, err
	}

	execTime := time.Now().UnixMilli()
	err = statement.QueryRow(id).Scan(&revenues.RevenueID, &revenues.MovieID, &revenues.DomesticTakings, &revenues.InternationalTakings)
	if err != nil {
		return revenues, err
	}
	fmt.Printf("Success! Execution time: %d milliseconds\n", time.Now().UnixMilli()-execTime)

	return revenues, err
}

// Sellect all rows from 'actors' table
func selectAllActors() ([]Actors, error) {
	var actors []Actors
	statement, err := db.Prepare("SELECT * FROM actors")
	if err != nil {
		return actors, err
	}

	execTime := time.Now().UnixMilli()
	rows, err := statement.Query()
	if err != nil {
		log.Fatal("Could not execute query: ", err)
	}

	for rows.Next() {
		var actor NillableActors
		if err := rows.Scan(&actor.ActorID, &actor.FirstName, &actor.LastName, &actor.Gender, &actor.DateOfBirth); err != nil {
			return actors, err
		}
		actors = append(actors, NillableToActors(actor))
	}
	fmt.Printf("Success! Execution time: %d milliseconds\n", time.Now().UnixMilli()-execTime)

	return actors, err
}

// Sellect all rows from 'directors' table
func selectAllDirectors() ([]Directors, error) {
	var directors []Directors
	statement, err := db.Prepare("SELECT * FROM directors")
	if err != nil {
		return directors, err
	}

	execTime := time.Now().UnixMilli()
	rows, err := statement.Query()
	if err != nil {
		log.Fatal("Could not execute query: ", err)
	}

	for rows.Next() {
		var director NillableDirectors
		if err := rows.Scan(&director.DirectorID, &director.FirstName, &director.LastName, &director.Nationality, &director.DateOfBirth); err != nil {
			return directors, err
		}
		directors = append(directors, NillableToDirectors(director))
	}
	fmt.Printf("Success! Execution time: %d milliseconds\n", time.Now().UnixMilli()-execTime)

	return directors, err
}

// Sellect all rows from 'movies' table
func selectAllMovies() ([]Movies, error) {
	var movies []Movies
	statement, err := db.Prepare("SELECT * FROM movies")
	if err != nil {
		return movies, err
	}

	execTime := time.Now().UnixMilli()
	rows, err := statement.Query()
	if err != nil {
		log.Fatal("Could not execute query: ", err)
	}

	for rows.Next() {
		var movie NillableMovies
		if err := rows.Scan(&movie.MovieID, &movie.MovieName, &movie.MovieLength, &movie.MovieLang, &movie.ReleaseDate, &movie.AgeCertificate, &movie.DirectorID); err != nil {
			return movies, err
		}
		movies = append(movies, NillableToMovies(movie))
	}
	fmt.Printf("Success! Execution time: %d milliseconds\n", time.Now().UnixMilli()-execTime)

	return movies, err
}

// Sellect all rows from 'movie_revenues' table
func selectAllMovieRevenues() ([]MovieRevenues, error) {
	var movieRevenues []MovieRevenues
	statement, err := db.Prepare("SELECT * FROM movie_revenues")
	if err != nil {
		return movieRevenues, err
	}

	execTime := time.Now().UnixMilli()
	rows, err := statement.Query()
	if err != nil {
		log.Fatal("Could not execute query: ", err)
	}

	for rows.Next() {
		var revenues NillableMovieRevenues
		if err := rows.Scan(&revenues.RevenueID, &revenues.MovieID, &revenues.DomesticTakings, &revenues.InternationalTakings); err != nil {
			return movieRevenues, err
		}
		movieRevenues = append(movieRevenues, NillableToMovieRevenues(revenues))
	}
	fmt.Printf("Success! Execution time: %d milliseconds\n", time.Now().UnixMilli()-execTime)

	return movieRevenues, err
}

//
// ---------- INSERT INTO SECTION ----------
//
// Insert a new row in 'actors' table and return its ID
func insertIntoActors(input Actors) int {
	var id int
	execTime := time.Now().UnixMilli()
	db.QueryRow("INSERT INTO actors (first_name, last_name, gender, date_of_birth) VALUES ($1, $2, $3, $4) RETURNING actor_id", input.FirstName, input.LastName, input.Gender, input.DateOfBirth).Scan(&id)

	fmt.Printf("Success! Execution time: %d milliseconds\n", time.Now().UnixMilli()-execTime)
	TablesLimit.ActorsLimit += 1
	fmt.Printf("Updating limit for table: actors\nLimit set - %d\n", TablesLimit.ActorsLimit)

	return id
}

// Insert a new row in 'directors' table and return its ID
func insertIntoDirectors(input Directors) int {
	var id int
	execTime := time.Now().UnixMilli()
	db.QueryRow("INSERT INTO directors (first_name, last_name, nationality, date_of_birth) VALUES ($1, $2, $3, $4) RETURNING director_id", input.FirstName, input.LastName, input.Nationality, input.DateOfBirth).Scan(&id)

	fmt.Printf("Success! Execution time: %d milliseconds\n", time.Now().UnixMilli()-execTime)
	TablesLimit.DirectorsLimit += 1
	fmt.Printf("Updating limit for table: directors\nLimit set - %d\n", TablesLimit.DirectorsLimit)

	return id
}

// Insert a new row in 'movies' table and return its ID
func insertIntoMovies(input Movies) int {
	var id int
	execTime := time.Now().UnixMilli()
	db.QueryRow("INSERT INTO movies (movie_name, movie_length, movie_lang, release_date, age_certificate, director_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING movie_id", input.MovieName, input.MovieLength, input.MovieLang, input.ReleaseDate, input.AgeCertificate).Scan(&id)

	fmt.Printf("Success! Execution time: %d milliseconds\n", time.Now().UnixMilli()-execTime)
	TablesLimit.MoviesLimit += 1
	fmt.Printf("Updating limit for table: movies\nLimit set - %d\n", TablesLimit.MoviesLimit)

	return id
}

// Insert a new row in 'directors' table and return its ID
func insertIntoMovieRevenues(input MovieRevenues) int {
	var id int
	execTime := time.Now().UnixMilli()
	db.QueryRow("INSERT INTO movie_revenues (revenue_id, movie_id, domestic_takings, international_takings) VALUES ($1, $2, $3, $4) RETURNING revenue_id", input.RevenueID, input.MovieID, input.DomesticTakings, input.InternationalTakings).Scan(&id)

	fmt.Printf("Success! Execution time: %d milliseconds\n", time.Now().UnixMilli()-execTime)
	TablesLimit.MovieRevenuesLimit += 1
	fmt.Printf("Updating limit for table: directors\nLimit set - %d\n", TablesLimit.MovieRevenuesLimit)

	return id
}

//
// ---------- DELETE FROM SECTION ----------
//
// Sellect a single row from 'actors' table using actor_id
func deleteFromActors(id int) (int, error) {
	var result int
	statement, err := db.Prepare("DELETE FROM actors WHERE actor_id = $1 RETURNING actor_id")
	if err != nil {
		return 0, err
	}

	execTime := time.Now().UnixMilli()
	err = statement.QueryRow(id).Scan(&result)
	if err != nil {
		return 0, err
	}
	fmt.Printf("Success! Execution time: %d milliseconds\n", time.Now().UnixMilli()-execTime)

	return result, err
}

// Sellect a single row from 'directors' table using actor_id
func deleteFromDirectors(id int) (int, error) {
	var result int
	statement, err := db.Prepare("DELETE FROM directors WHERE director_id = $1 RETURNING director_id")
	if err != nil {
		return 0, err
	}

	execTime := time.Now().UnixMilli()
	err = statement.QueryRow(id).Scan(&result)
	if err != nil {
		return 0, err
	}
	fmt.Printf("Success! Execution time: %d milliseconds\n", time.Now().UnixMilli()-execTime)

	return result, err
}

// Sellect a single row from 'movies' table using actor_id
func deleteFromMovies(id int) (int, error) {
	var result int
	statement, err := db.Prepare("DELETE FROM movies WHERE movie_id = $1 RETURNING movie_id")
	if err != nil {
		return 0, err
	}

	execTime := time.Now().UnixMilli()
	err = statement.QueryRow(id).Scan(&result)
	if err != nil {
		return 0, err
	}
	fmt.Printf("Success! Execution time: %d milliseconds\n", time.Now().UnixMilli()-execTime)

	return result, err
}

// Sellect a single row from 'movie_revenues' table using actor_id
func deleteFromMovieRevenues(id int) (int, error) {
	var result int
	statement, err := db.Prepare("DELETE FROM movie_revenues WHERE revenue_id = $1 RETURNING revenue_id")
	if err != nil {
		return 0, err
	}

	execTime := time.Now().UnixMilli()
	err = statement.QueryRow(id).Scan(&result)
	if err != nil {
		return 0, err
	}
	fmt.Printf("Success! Execution time: %d milliseconds\n", time.Now().UnixMilli()-execTime)

	return result, err
}

*/
