package handlers

import (
	"movies-api/internal/repository"
	"movies-api/internal/service"
)

// App stores the services for different entities
type App struct {
	Repo         *repository.Repo // to be deleted once all services are set
	MovieService *service.MovieService
}
