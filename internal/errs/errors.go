package errs

import (
	"context"
	"errors"
	"log"
	"movies-api/internal/repository"
	"net/http"
)

// Invalid user input
var ErrInvalidUserInput = errors.New("invalid input")

var ErrForce = errors.New("operation requires force")

// WriteError chooses the appropriate error response and writes to http.ResponseWrite
func WriteError(w http.ResponseWriter, err error) {
	if errors.Is(err, context.Canceled) {
		log.Println("user disconnected before response finished")
		return
	}

	if errors.Is(err, repository.ErrNotFound) {
		http.Error(w, repository.ErrNotFound.Error(), http.StatusNotFound)
		return
	}

	if errors.Is(err, ErrInvalidUserInput) {
		http.Error(w, err.Error(), http.StatusBadRequest) // Show detailed error so user can fix input
		return
	}

	if errors.Is(err, ErrForce) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
