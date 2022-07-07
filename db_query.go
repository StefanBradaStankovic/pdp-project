package main

/*
//
// ---------------------------------------- INSERT INTO SECTION ----------------------------------------
// ---------------------------------------- INSERT INTO SECTION ----------------------------------------
// ---------------------------------------- INSERT INTO SECTION ----------------------------------------
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
}*/
