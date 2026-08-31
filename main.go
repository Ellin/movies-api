package main

import (
	"database/sql"
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

	validate := initValidator()

	app := initApp(db, validate)

	if err := setupDB(db, app, reset); err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Database set up successfully")

	// start server
	fmt.Println("Starting server...")
	log.Fatalln(http.ListenAndServe(":8080", NewRouter(app)))
}

// parseFlags parses the CLI flags, returning the database file name and reset bool
func parseFlags() (string, bool) {
	dbFile := flag.String("db", "movies.db", "Database file")
	reset := flag.Bool("reset", false, "Resets database and seeds it with dummy data (default false)")

	flag.Parse()
	return *dbFile, *reset
}

// initValidator creates a new *validator.Validate and registers the JSON field names
func initValidator() *validator.Validate {
	validate := validator.New()

	// Register JSON field names as validator field names instead of default Go struct field names
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return validate
}

// initApp initialize the application's dependencies
func initApp(db *sql.DB, validate *validator.Validate) *handlers.App {
	repo := &repository.Repo{DB: db}

	app := handlers.App{
		MovieService: service.NewMovieService(repo, validate),
		ActorService: service.NewActorService(repo, validate),
		GenreService: service.NewGenreService(repo, validate),
	}

	return &app
}

// setupDB sets up the database, reseeding it if reset flag is true
func setupDB(db *sql.DB, app *handlers.App, reset bool) error {
	if reset {
		// Clear database and seed with dummy data
		if err := database.ResetDatabase(db, app); err != nil {
			return fmt.Errorf("reset database: %w", err)
		}
	} else {
		// Create database tables without any dummy data
		if err := database.CreateTables(db); err != nil {
			return fmt.Errorf("create database tables: %w", err)
		}
	}

	return nil
}
