package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

//
// -------------------- QUERY SECTION --------------------
// -------------------- QUERY SECTION --------------------
// -------------------- QUERY SECTION --------------------
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

//
// -------------------- CONTROLL SECTION --------------------
// -------------------- CONTROLL SECTION --------------------
// -------------------- CONTROLL SECTION --------------------
//
// Make a DB connection string using specified config
func setDBConnection() string {
	psqlConnect := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", db_host, db_port, db_user, db_password, db_name)

	return psqlConnect
}

// Connect to the database
func dbConnect() *sql.DB {
	fmt.Printf("Connecting to database...	")
	db, err := sql.Open("postgres", psqlConnect)
	if err != nil {
		log.Fatal("ERROR - db_controll.go - Could not connect to the database: ", err)
	} else {
		fmt.Printf("Success!\n")
	}

	return db
}

// Prepare a database statement
func checkIfExists(id int, queryString string) (bool, error) {
	statement, err := db.Prepare(queryString)
	var result int
	if err != nil {
		fmt.Printf("ERROR - db_controll.go - %s\n", err)
		return false, err
	}
	err = statement.QueryRow(id).Scan(&result)
	fmt.Printf("Query result is: %d\n", result)
	if err != nil || result <= 0 {
		return false, err
	}

	return true, err
}
