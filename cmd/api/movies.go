package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Ramdoni007/21Cinema/internal/data"
	"github.com/Ramdoni007/21Cinema/internal/validator"
)

// Add a createMovieHandler for the "POST /v1/movies" endpoint. For now we simply
// return a plain-text placeholder response.
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	// Use the new readJSON() helper to decode the request body into the input struct.
	// If this returns an error we send the client the error message along with a 400
	// Bad Request status code, just like before.
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	// Initialize a new Validator instance.
	v := validator.New()

	if data.ValidateMovie(v, movie); !v.Valid() {
		// Use the new failedValidationResponse () helper
		app.failedValidationResponse(w, r, v.Errors)
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
		Title:     "BlackHawkDown",
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
