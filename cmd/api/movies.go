package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Ramdoni007/21Cinema/internal/data"
)

// Add a createMovieHandler for the "POST /v1/movies" endpoint. For now we simply
// return a plain-text placeholder response.
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Tittle  string   `json:"tittle"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Fprintf(w, "%+v\n", input)
}

// Add a showMovieHandler for the "GET /v1/movies/:id" endpoint. For now, we retrieve
// the interpolated "id" parameter from the current URL and include it in a placeholder
// response.
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	// Use Method readIDParams(r) to extract "id" Parameter from URL
	id, err := app.readIDParams(r)
	if err != nil {
		// Use the new notFoundResponse() helper
		app.notFoundResponse(w, r)
		return
	}

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Tittle:    "BlackHawkDown",
		Year:      2003,
		Runtime:   102,
		Genres:    []string{"drama", "conflict", "war", "army"},
		Version:   1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		// Use The new serverErrorResponse() helper
		app.serverErrorResponse(w, r, err)
	}
}
