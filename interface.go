package main

import (
	"database/sql"
	"fmt"
)

// Wrapper for database/sql type *sql.Row since golang interfaces
// do not support that type as receiver
type SqlRow struct {
	row *sql.Row
}

// Interface for scanning *sql.Row elements into struct objects
type rowScanner interface {
	ScanActor() Actors
	ScanDirector() Directors
}

// rowScanner method for scanning an 'Actors' object
func (inputRow *SqlRow) ScanActor() Actors {
	var output Actors
	err := inputRow.row.Scan(&output.ActorID, &output.FirstName, &output.LastName, &output.Gender, &output.DateOfBirth)
	if err != nil {
		fmt.Printf("ERROR - interface.go -  %s\n", err)
		return output
	}
	return output
}

// rowScanner method for scanning a 'Directors' object
func (inputRow *SqlRow) ScanDirector() Directors {
	var output Directors
	err := inputRow.row.Scan(&output.DirectorID, &output.FirstName, &output.LastName, &output.Nationality, &output.DateOfBirth)
	if err != nil {
		fmt.Printf("ERROR - interface.go -  %s\n", err)
		return output
	}
	return output
}

// rowScanner method for scanning a 'Movies' object
func (inputRow *SqlRow) ScanMovie() Movies {
	var output Movies
	err := inputRow.row.Scan(&output.MovieID, &output.MovieName, &output.MovieLength, &output.MovieLang, &output.ReleaseDate, &output.AgeCertificate, &output.DirectorID)
	if err != nil {
		fmt.Printf("ERROR - interface.go -  %s\n", err)
		return output
	}
	return output
}

// rowScanner method for scanning a 'MovieRevenues' object
func (inputRow *SqlRow) ScanRevenues() MovieRevenues {
	var output MovieRevenues
	err := inputRow.row.Scan(&output.RevenueID, &output.MovieID, &output.DomesticTakings, &output.InternationalTakings)
	if err != nil {
		fmt.Printf("ERROR - interface.go -  %s\n", err)
		return output
	}
	return output
}
