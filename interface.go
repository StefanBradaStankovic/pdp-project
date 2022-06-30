package main

import "database/sql"

// DB configuration
const (
	db_host     = "localhost"
	db_port     = 5432
	db_user     = "postgres"
	db_password = "htec1234"
	db_name     = "movies_data"
)

// DB table row count
type DBTableLimit struct {
	ActorsLimit        int
	DirectorsLimit     int
	MovieRevenuesLimit int
	MoviesLimit        int
}

type Actors struct {
	ActorID     int    `json:"actorID"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Gender      string `json:"gender"`
	DateOfBirth string `json:"dateOfBirth"`
}

type Directors struct {
	DirectorID  int    `json:"directorID"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Nationality string `json:"nationality"`
	DateOfBirth string `json:"dateOfBirth"`
}

type Movies struct {
	MovieID        int    `json:"movieID"`
	MovieName      string `json:"movieName"`
	MovieLength    string `json:"movieLength"`
	MovieLang      string `json:"movieLang"`
	ReleaseDate    string `json:"releaseDate"`
	AgeCertificate string `json:"ageCertificate"`
	DirectorID     int    `json:"directorID"`
}

type MovieRevenues struct {
	RevenueID            int     `json:"revenueID"`
	MovieID              int     `json:"movieID"`
	DomesticTakings      float64 `json:"domesticTakings"`
	InternationalTakings float64 `json:"internationalTakings"`
}

// NILLABLE STRUCTURES
type NillableActors struct {
	ActorID     sql.NullInt64  `json:"actorID"`
	FirstName   sql.NullString `json:"firstName"`
	LastName    sql.NullString `json:"lastName"`
	Gender      sql.NullString `json:"gender"`
	DateOfBirth sql.NullString `json:"dateOfBirth"`
}

type NillableDirectors struct {
	DirectorID  sql.NullInt64  `json:"directorID"`
	FirstName   sql.NullString `json:"firstName"`
	LastName    sql.NullString `json:"lastName"`
	Nationality sql.NullString `json:"nationality"`
	DateOfBirth sql.NullString `json:"dateOfBirth"`
}

type NillableMovies struct {
	MovieID        sql.NullInt64  `json:"movieID"`
	MovieName      sql.NullString `json:"movieName"`
	MovieLength    sql.NullString `json:"movieLength"`
	MovieLang      sql.NullString `json:"movieLang"`
	ReleaseDate    sql.NullString `json:"releaseDate"`
	AgeCertificate sql.NullString `json:"ageCertificate"`
	DirectorID     sql.NullInt64  `json:"directorID"`
}

type NillableMovieRevenues struct {
	RevenueID            sql.NullInt64   `json:"revenueID"`
	MovieID              sql.NullInt64   `json:"movieID"`
	DomesticTakings      sql.NullFloat64 `json:"domesticTakings"`
	InternationalTakings sql.NullFloat64 `json:"internationalTakings"`
}

// CONVERSION FROM NILLABLE TO REGULAR
func NillableToActors(input NillableActors) Actors {
	var output Actors

	output.ActorID = int(input.ActorID.Int64)
	output.FirstName = input.FirstName.String
	output.LastName = input.LastName.String
	output.Gender = input.Gender.String
	output.DateOfBirth = input.DateOfBirth.String

	return output
}

func NillableToDirectors(input NillableDirectors) Directors {
	var output Directors

	output.DirectorID = int(input.DirectorID.Int64)
	output.FirstName = input.FirstName.String
	output.LastName = input.LastName.String
	output.Nationality = input.Nationality.String
	output.DateOfBirth = input.DateOfBirth.String

	return output
}

func NillableToMovies(input NillableMovies) Movies {
	var output Movies

	output.MovieID = int(input.DirectorID.Int64)
	output.MovieName = input.MovieName.String
	output.MovieLength = input.MovieLength.String
	output.MovieLang = input.MovieLang.String
	output.ReleaseDate = input.ReleaseDate.String
	output.AgeCertificate = input.AgeCertificate.String
	output.DirectorID = int(input.DirectorID.Int64)

	return output
}

func NillableToMovieRevenues(input NillableMovieRevenues) MovieRevenues {
	var output MovieRevenues

	output.RevenueID = int(input.RevenueID.Int64)
	output.MovieID = int(input.MovieID.Int64)
	output.DomesticTakings = input.DomesticTakings.Float64
	output.InternationalTakings = input.InternationalTakings.Float64

	return output
}
