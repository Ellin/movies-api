package handlers

import (
	"movies-api/internal/service"
)

// App stores the services for different entities
type App struct {
	MovieService *service.MovieService
	GenreService *service.GenreService
	ActorService *service.ActorService
}
