package handlers

import (
	"encoding/json"
	"fmt"
	"movies-api/internal/errs"
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
		errs.WriteError(w, err)
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
		errs.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(actors)
}

// get
func (app *App) GetActor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		errs.WriteError(w, fmt.Errorf("%w: %w", errs.ErrInvalidUserInput, err))
		return
	}

	actor, err := app.ActorService.GetActor(ctx, id)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(actor)
}

func (app *App) PatchActor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		errs.WriteError(w, fmt.Errorf("%w: %w", errs.ErrInvalidUserInput, err))
		return
	}

	patch := service.ActorPatch{}

	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	actor, err := app.ActorService.PatchActor(ctx, id, patch)
	if err != nil {
		errs.WriteError(w, err)
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
		errs.WriteError(w, fmt.Errorf("%w: %w", errs.ErrInvalidUserInput, err))
		return
	}
	var force bool = r.URL.Query().Get("force") == "true"
	if err := app.ActorService.DeleteActor(ctx, id, force); err != nil {
		errs.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
