package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// Declare an envelope map containing the data for the response. Notice that the way
	// we've constructed this means the environment and version data will now be nested
	// under a system_info key in the JSON response.
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	// writeJSON Creating a writeJSON helper method with helpers.go
	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(
			w,
			"the server encountered a problem and could not process your request",
			http.StatusInternalServerError,
		)

	}
}
