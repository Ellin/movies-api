package main

import "database/sql"

func CreateTableMovies(db *sql.DB) (sql.Result, error) {
	queryMovies := `CREATE TABLE IF NOT EXISTS movies (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL,
	releaseYear INTEGER NOT NULL,
	duration INTEGER NOT NULL
	)`

	return db.Exec(queryMovies)
}

func CreateTableActors(db *sql.DB) (sql.Result, error) {
	queryActors := `CREATE TABLE IF NOT EXISTS actors (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	birthDate TEXT NOT NULL
	)`

	return db.Exec(queryActors)
}

func CreateTableGenres(db *sql.DB) (sql.Result, error) {
	queryGenres := `CREATE TABLE IF NOT EXISTS genres (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL
	)`

	return db.Exec(queryGenres)
}

func JoinTableMoviesGenres(db *sql.DB) (sql.Result, error) {
	queryMoviesGenres := `CREATE TABLE IF NOT EXISTS movies_genres (
	movie_id INTEGER NOT NULL,
	genre_id INTEGER NOT NULL,
	FOREIGN KEY movie_id REFERENCES movies(id),
	FOREIGN KEY genre_id REFERENCES movies(id),
	)`

	return db.Exec(queryMoviesGenres)
}
