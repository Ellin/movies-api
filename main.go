package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./movies.db")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()

	fmt.Println("Connected to the SQLite database successfully.")

	// var sqliteversion string
	// err = database.QueryRow("select sqlite_version()").Scan(&sqliteversion)

	// fmt.Println(sqliteversion)

	//Create table for Movies
	_, err = CreateTableMovies(db)
	if err != nil {
		fmt.Println(err)
		return
	}

	//Create table for Actors
	_, err = CreateTableActors(db)
	if err != nil {
		fmt.Println(err)
		return
	}

	//Create table for Genres
	_, err = CreateTableGenres(db)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("All tables created successfully")
}
