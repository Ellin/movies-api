package main

import (
	"flag"
	"fmt"
	"log"
	"movies-api/internal/database"
	"movies-api/internal/handlers"
	"movies-api/internal/repository"
	"movies-api/internal/service"
	"net/http"

	"github.com/go-playground/validator/v10"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	dbFile, reset := parseFlags()

	db, err := database.OpenDB(dbFile + "?_foreign_keys=on") // enforce foreign keys -> validate existence of rows referred to by foreign keys
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
	validate := validator.New()
	app := handlers.App{
		Repo:         repo, // to be deleted once all services are set
		MovieService: service.NewMovieService(repo, validate),
		GenreService: service.NewGenreService(repo),
		ActorService: service.NewActorService(repo),
	}

	// Clear database data and seed with dummy data
	if reset {
		err = database.ResetDatabase(&app)
		if err != nil {
			log.Fatalln("Reset database failure:", err)
		}
	}

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
	// _, err = app.Repo.AddActor(ctx, models.Actor{Name: "Steven King", BirthDate: "25.05.1967"})
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// actor := models.ActorDetail{}
	// fmt.Println("Finished adding actor")
	// actor, err = repo.GetActor(ctx, 1)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(actor.Name)

	// start server
	fmt.Println("Starting server...")
	log.Fatalln(http.ListenAndServe(":8080", NewRouter(&app)))
}

// parseFlags parses the CLI flags, returning the database file name and reset bool
func parseFlags() (string, bool) {
	dbFile := flag.String("db", "movies.db", "Database file")
	reset := flag.Bool("reset", false, "Resets database and seeds it with dummy data (default false)")

	flag.Parse()

	return *dbFile, *reset
}
