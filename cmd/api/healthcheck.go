package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// Update Code Sending JSON response to EndPoint healthcheck
	// Create a map which holds the information that we want to send in the response.
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}

	// writeJSON Creating a writeJSON helper method with helpers.go
	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(
			w,
			"the server encountered a problem and could not process your request",
			http.StatusInternalServerError,
		)

	}
}
