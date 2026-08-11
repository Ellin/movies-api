package main

import (
	"fmt"
	"log"
	"movies-api/internal/database"
	"movies-api/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

// application stores the app's dependencies
type application struct {
	repo *models.Repo
}

func main() {
	db, err := database.OpenDB("./movies.db?_foreign_keys=on") // enforce foreign keys -> validate existence of rows referred to by foreign keys
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	fmt.Println("Connected to SQLite database successfully.")

	err = database.InitDB(db)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("Database tables initialized successfully.")

	// Initialize the application's dependencies
	app := application{
		repo: &models.Repo{DB: db},
	}
	fmt.Println(app) // temp use of app so go doesn't give error
}
