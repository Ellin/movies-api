package main

import (
	"context"
	"fmt"
	"log"
	"movies-api/internal/database"
	"movies-api/internal/handlers"
	"movies-api/internal/models"
	"movies-api/internal/repository"
	"movies-api/internal/service"

	_ "github.com/mattn/go-sqlite3"
)

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
	app := handlers.App{
		Repo:         repo, // to be deleted once all services are set
		ActorService: service.NewActorService(repo),
	}
	ctx := context.TODO()

	// genr := models.Genre{ID: 1, Name: "Drama"}

	// res, err := app.repo.CreateGenre(genr)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println("Created with:", res)

	// fmt.Println(app.repo.GetGenreByID(res))

	// fmt.Println(app.repo.GetAllGenres())

	// //movies
	// _, err = app.repo.AddMovie(nil, models.Movie{Title: "Spiderman", ReleaseYear: 2000, Duration: 120})
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// //
	//actors
	_, err = app.Repo.AddActor(ctx, models.Actor{Name: "Steven King", BirthDate: "25.05.1967"})
	if err != nil {
		fmt.Println(err)
	}
	actor := models.ActorSummary{}
	fmt.Println("Finished adding actor")
	actor, err = repo.GetActor(ctx, 1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(actor.Name)

	// start server
	// fmt.Println("Starting server...")
	// err = http.ListenAndServe(":8080", NewRouter(&app))
	// if err != nil {
	// 	log.Fatalln("starting server:", err)
	// }

}
