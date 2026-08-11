package main

import (
	"fmt"
	"log"
	"movies-api/internal/database"
	"movies-api/internal/models"
	"movies-api/internal/repository"

	_ "github.com/mattn/go-sqlite3"
)

// application stores the app's dependencies
type application struct {
	repo *repository.Repo
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
		repo: &repository.Repo{DB: db},
	}

	genr := models.Genre{ID: 1, Name: "blabla"}

	res, err := app.repo.CreateGenre(genr)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Created with:", res)

	fmt.Println(app.repo.GetGenreByID(res))

}
