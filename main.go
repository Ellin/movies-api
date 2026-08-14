package main

import (
	"fmt"
	"log"
	"movies-api/internal/database"
	"movies-api/internal/models"
	"movies-api/internal/repository"
	"movies-api/internal/service"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// application stores the app's dependencies
type application struct {
	repo         *repository.Repo
	movieService *service.MovieService
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
	repo := &repository.Repo{DB: db}
	app := application{
		repo:         repo,
		movieService: service.NewMovieService(repo),
	}

	genr := models.Genre{ID: 1, Name: "Drama"}

	res, err := app.repo.CreateGenre(genr)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Created with:", res)

	fmt.Println(app.repo.GetGenreByID(res))

	fmt.Println(app.repo.GetAllGenres())

	// movies
	_, err = app.repo.AddMovie(models.Movie{Title: "Bambie", ReleaseYear: 2023, Duration: 90})
	if err != nil {
		fmt.Println(err)
	}

	// start server
	fmt.Println("Starting server...")
	err = http.ListenAndServe(":8080", app.routes())
	if err != nil {
		log.Fatalln("starting server:", err)
	}

}
