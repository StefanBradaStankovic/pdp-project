package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var TablesLimit DBTableLimit
var psqlConnect string
var db *sql.DB

//
// ---------- CONTROLL SECTION ----------
//
// Make a DB connection string using specified
func setDBConnection() string {
	psqlConnect := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", db_host, db_port, db_user, db_password, db_name)

	return psqlConnect
}

// Connect to the database
func dbConnect() *sql.DB {
	fmt.Printf("Connecting to database...	")
	db, err := sql.Open("postgres", psqlConnect)
	if err != nil {
		log.Fatal("Could not connect to the database: ", err)
	} else {
		fmt.Printf("Success!\n")
	}

	return db
}

// Count DB rows in each table
func setTablesLimit() {
	var err error

	fmt.Printf("Setting limit for table: actors\n")
	err = db.QueryRow("SELECT MAX(actor_id) FROM actors").Scan(&TablesLimit.ActorsLimit)
	if err != nil {
		log.Fatal("Could not execute query: ", err)
	} else {
		fmt.Printf("Limit set - maximum ID detected '%d'\n", TablesLimit.ActorsLimit)
	}

	fmt.Printf("Setting limit for table: directors\n")
	err = db.QueryRow("SELECT MAX(director_id) FROM directors").Scan(&TablesLimit.DirectorsLimit)
	if err != nil {
		log.Fatal("Could not execute query: ", err)
	} else {
		fmt.Printf("Limit set - maximum ID detected '%d'\n", TablesLimit.DirectorsLimit)
	}

	fmt.Printf("Setting limit for table: movies\n")
	err = db.QueryRow("SELECT MAX(movie_id) FROM movies").Scan(&TablesLimit.MoviesLimit)
	if err != nil {
		log.Fatal("Could not execute query: ", err)
	} else {
		fmt.Printf("Limit set - maximum ID detected '%d'\n", TablesLimit.MoviesLimit)
	}

	fmt.Printf("Setting limit for table: movie_revenues\n")
	err = db.QueryRow("SELECT MAX(movie_id) FROM movie_revenues").Scan(&TablesLimit.MovieRevenuesLimit)
	if err != nil {
		log.Fatal("Could not execute query: ", err)
	} else {
		fmt.Printf("Limit set - maximum ID detected '%d'\n", TablesLimit.MovieRevenuesLimit)
	}

}
