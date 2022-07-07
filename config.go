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
var RowScanner rowScanner
var psqlConnect string
var db *sql.DB

type Directors struct {
	DirectorID  *int    `json:"directorID"`
	FirstName   *string `json:"firstName"`
	LastName    *string `json:"lastName"`
	Nationality *string `json:"nationality"`
	DateOfBirth *string `json:"dateOfBirth"`
}

type Movies struct {
	MovieID        *int    `json:"movieID"`
	MovieName      *string `json:"movieName"`
	MovieLength    *string `json:"movieLength"`
	MovieLang      *string `json:"movieLang"`
	ReleaseDate    *string `json:"releaseDate"`
	AgeCertificate *string `json:"ageCertificate"`
	DirectorID     *int    `json:"directorID"`
}

type MovieRevenues struct {
	RevenueID            *int     `json:"revenueID"`
	MovieID              *int     `json:"movieID"`
	DomesticTakings      *float64 `json:"domesticTakings"`
	InternationalTakings *float64 `json:"internationalTakings"`
}
