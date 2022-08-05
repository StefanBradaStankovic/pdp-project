package main

import (
	"database/sql"
)

// DB connection configuration
const (
	db_host     = "localhost"
	db_port     = 5432
	db_user     = "postgres"
	db_password = "htec1234"
	db_name     = "movies_data"
)

// Global variables
//var RowScanner rowScanner
var psqlConnect string
var db *sql.DB

type queryStatement struct {
	selectSingle  string
	selectAll     string
	insertInto    string
	deleteFrom    string
	updateVisible string
	checkForID    string
}

// Wrapper for database/sql type *sql.Row since golang interfaces
// do not support that type as receiver
type SqlRow struct {
	row *sql.Row
}
type SqlRows struct {
	rows *sql.Rows
}
type SqlStmt struct {
	statement *sql.Stmt
}
