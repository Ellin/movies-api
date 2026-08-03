package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func CreateTable(db *sql.DB) (sql.Result, error) {
	query := `CREATE TABLE IF NOT EXISTS movies (
	id INTEGER PRIMARY KEY,
	title TEXT NOT NULL,
	releaseYear INTEGER NOT NULL,
	duration INTEGER NOT NULL
	)`

	return db.Exec(query)
}

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
	//
	_, err = CreateTable(db)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Table created successfully")
}
