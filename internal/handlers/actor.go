package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"movies-api/internal/repository"
	"movies-api/internal/service"
	"net/http"
	"strconv"
)

func (app *App) PostActor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sub := service.ActorSubmission{}

	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	actor, err := app.ActorService.AddActor(ctx, sub)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			log.Println("client disconnected before add movie finished")
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(actor)
}

// get
func (app *App) GetAllActors(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	actors, err := app.ActorService.GetAllActors(ctx)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			log.Println("client disconnected before get actor finished")
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actors)
}

// get
func (app *App) GetActor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	actors, err := app.ActorService.GetActor(ctx, id)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			log.Println("client disconnected before get movie finished")
		} else if errors.Is(err, repository.ErrNotFound) {
			http.NotFound(w, r)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actors)
}

func (app *App) PatchActor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	patch := service.ActorPatch{}

	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	actor, err := app.ActorService.PatchActor(ctx, id, patch)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			log.Println("client disconnected before add movie finished")
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)

			// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(actor)
}

func (app *App) DeleteActor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	if err := app.ActorService.DeleteActor(ctx, id); err != nil {
		if errors.Is(err, context.Canceled) {
			log.Println("client disconnected before get movie finished")
		} else if errors.Is(err, repository.ErrNotFound) {
			http.NotFound(w, r)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
