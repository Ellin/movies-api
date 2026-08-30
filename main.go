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
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Get CLI flags
	dbFile, reset := parseFlags()

	// Initialize database connection pool
	db, err := database.OpenDB(dbFile + "?_foreign_keys=on") // enforce foreign keys -> validate existence of rows referred to by foreign keys
	if err != nil {
		log.Fatalln("Database connection failure:", err)
	}
	defer db.Close()
	fmt.Printf("Connected to SQLite database %s\n", dbFile)

	// Initialize the application's dependencies
	repo := &repository.Repo{DB: db}
	// Intialize validator
	validate := validator.New()

	// Register JSON field names as validator field names instead of default Go struct field names
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	app := handlers.App{
		Repo:         repo, // to be deleted once all services are set
		MovieService: service.NewMovieService(repo, validate),
		ActorService: service.NewActorService(repo, validate),
		GenreService: service.NewGenreService(repo, validate),
	}

	// Check database reset flag
	if reset {
		// Clear database and seed with dummy data
		if err := database.ResetDatabase(&app); err != nil {
			log.Fatalln("Reset database failure:", err)
			return
		}
	} else {
		// Initialize database tables without any dummy data
		if err := database.InitDB(db); err != nil {
			log.Fatalln("Database initialization failure:", err)
			return
		}
	}
	fmt.Println("Database initialized successfully")

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
