package main

import (
	"fmt"
	"log"
	"movies-api/internal/database"
	"movies-api/internal/repository"

	_ "github.com/mattn/go-sqlite3"
)

<<<<<<< HEAD
=======
// application stores the app's dependencies
type application struct {
	repo *repository.Repo
}

>>>>>>> acd27f2cebc25a46127a7966e9b1220b089e4668
func main() {
	db, err := database.OpenDB("./movies.db?_foreign_keys=on") // enforce foreign keys -> validate existence of rows referred to by foreign keys
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	fmt.Println("Connected to SQLite database successfully.")

<<<<<<< HEAD
	// var sqliteversion string
	// err = database.QueryRow("select sqlite_version()").Scan(&sqliteversion)

	// fmt.Println(sqliteversion)

	//Create table for Movies
	_, err = CreateTableMovies(db)
=======
	err = database.InitDB(db)
>>>>>>> acd27f2cebc25a46127a7966e9b1220b089e4668
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("Database tables initialized successfully.")

<<<<<<< HEAD
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
=======
	// Initialize the application's dependencies
	app := application{
		repo: &repository.Repo{DB: db},
	}
	fmt.Println(app) // temp use of app so go doesn't give error
>>>>>>> acd27f2cebc25a46127a7966e9b1220b089e4668
}
