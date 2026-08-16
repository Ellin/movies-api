package database

import (
	"database/sql"
	"fmt"
)

// OpenDB initializes a database connection pool for a given data source name (dns) or connection string.
// The connection is tested before returning the sql.DB connection pool.
func OpenDB(dsn string) (*sql.DB, error) {
	// Initialize database connection pool
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	// Create database connection and check for errors
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func createTableGenres(db *sql.DB) (sql.Result, error) {
	sql := `CREATE TABLE IF NOT EXISTS genres (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL
	);`

	return db.Exec(sql)
}

func createTableMovies(db *sql.DB) (sql.Result, error) {
	sql := `CREATE TABLE IF NOT EXISTS movies (
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL,
		releaseYear INTEGER NOT NULL,
		duration INTEGER NOT NULL
	);`

	return db.Exec(sql)
}

func createTableActors(db *sql.DB) (sql.Result, error) {
	sql := `CREATE TABLE IF NOT EXISTS actors (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		birthDate TEXT NOT NULL
	);`

	return db.Exec(sql)
}

func createTableGenresMovies(db *sql.DB) (sql.Result, error) {
	sql := `CREATE TABLE IF NOT EXISTS genres_movies (
		genre_id INTEGER,
		movie_id INTEGER,
		PRIMARY KEY (genre_id, movie_id),
		FOREIGN KEY (genre_id) REFERENCES genres (id),
		FOREIGN KEY (movie_id) REFERENCES movies (id) ON DELETE CASCADE
	);`

	return db.Exec(sql)
}

func createTableMoviesActors(db *sql.DB) (sql.Result, error) {
	sql := `CREATE TABLE IF NOT EXISTS movies_actors (
		movie_id INTEGER,
		actor_id INTEGER,
		PRIMARY KEY (movie_id, actor_id),
		FOREIGN KEY (movie_id) REFERENCES movies (id) ON DELETE CASCADE,
		FOREIGN KEY (actor_id) REFERENCES actors (id)
	);`

	return db.Exec(sql)
}

func InitDB(db *sql.DB) error {
	_, err := createTableGenres(db)
	if err != nil {
		return fmt.Errorf("Creating genres table: %w", err)
	}

	_, err = createTableMovies(db)
	if err != nil {
		return fmt.Errorf("Creating movies table: %w", err)
	}

	_, err = createTableActors(db)
	if err != nil {
		return fmt.Errorf("Creating actors table: %w", err)
	}

	_, err = createTableGenresMovies(db)
	if err != nil {
		return fmt.Errorf("Creating genres_movies table: %w", err)
	}

	_, err = createTableMoviesActors(db)
	if err != nil {
		return fmt.Errorf("Creating movies_actors table: %w", err)
	}

	return nil
}
