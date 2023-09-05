package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Ramdoni007/21Cinema/internal/data"
)

// Add a createMovieHandler for the "POST /v1/movies" endpoint. For now we simply
// return a plain-text placeholder response.
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new movie")
}

// Add a showMovieHandler for the "GET /v1/movies/:id" endpoint. For now, we retrieve
// the interpolated "id" parameter from the current URL and include it in a placeholder
// response.
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	// Use Method readIDParams(r) to extract "id" Parameter from URL
	id, err := app.readIDParams(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Tittle:    "BlackHawkDown",
		Year:      2003,
		Runtime:   007,
		Genres:    []string{"drama", "conflict", "war", "army"},
		Version:   1,
	}

	err = app.writeJSON(w, http.StatusOK, movie, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(
			w,
			"the server encountered a problem and could not process your request",
			http.StatusInternalServerError,
		)

	}
}
